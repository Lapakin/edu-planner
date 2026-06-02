package generator

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"
)

// flowKey identifies a flow stream: lessons of the same subject taught by the same
// teacher are candidates to be merged into a shared lecture for several groups.
type flowKey struct {
	subjectID uint64
	teacherID uint64
}

// mergeFlowLessons merges flow-eligible lessons (a flow discipline's lecturer
// lessons sharing the same teacher across several groups) into shared flow
// lessons. Lessons that have no flow partner stay as regular single-group lessons.
//
// Merging is independent per call (one week period). When the per-group lesson
// counts differ, the surplus lessons of a group simply remain separate, which is a
// natural partial fallback.
func mergeFlowLessons(lessons []*lesson, cfg *settings) []*lesson {
	if cfg != nil && !cfg.allowFlowLessons {
		return lessons
	}

	var others, candidates []*lesson
	for _, l := range lessons {
		if isFlowCandidate(l) {
			candidates = append(candidates, l)
		} else {
			others = append(others, l)
		}
	}
	if len(candidates) == 0 {
		return lessons
	}

	// Group candidates by (subject, teacher), then by group within each stream.
	groupsByKey := make(map[flowKey]map[uint64][]*lesson)
	keyOrder := make([]flowKey, 0)
	for _, l := range candidates {
		k := flowKey{subjectID: l.subjectID, teacherID: l.subLessons[0].teacherID}
		if groupsByKey[k] == nil {
			groupsByKey[k] = make(map[uint64][]*lesson)
			keyOrder = append(keyOrder, k)
		}
		groupsByKey[k][l.groupID] = append(groupsByKey[k][l.groupID], l)
	}
	sort.Slice(keyOrder, func(i, j int) bool {
		if keyOrder[i].subjectID != keyOrder[j].subjectID {
			return keyOrder[i].subjectID < keyOrder[j].subjectID
		}
		return keyOrder[i].teacherID < keyOrder[j].teacherID
	})

	result := others
	for _, k := range keyOrder {
		byGroup := groupsByKey[k]
		gids := sortedKeys(byGroup)

		for {
			// Groups that still have an unconsumed lesson this round.
			active := make([]uint64, 0, len(gids))
			for _, gid := range gids {
				if len(byGroup[gid]) > 0 {
					active = append(active, gid)
				}
			}

			// Fewer than two groups left → no flow possible; emit leftovers as
			// regular single-group lessons (they were never promoted to flow).
			if len(active) < 2 {
				for _, gid := range active {
					result = append(result, byGroup[gid]...)
					byGroup[gid] = nil
				}
				break
			}

			// Take one lesson from each active group and merge them.
			constituents := make([]*lesson, 0, len(active))
			for _, gid := range active {
				constituents = append(constituents, byGroup[gid][0])
				byGroup[gid] = byGroup[gid][1:]
			}
			result = append(result, mergeFlow(constituents, k))
		}
	}

	return result
}

// isFlowCandidate reports whether a lesson may participate in a flow: a
// flow-eligible, non-lab, whole-group (united) lecture with at least one
// sub-lesson.
func isFlowCandidate(l *lesson) bool {
	return l.flowEligible && !l.isLab && l.format == formatUnited && len(l.subLessons) > 0
}

// mergeFlow combines per-group constituent lessons into a single flow lesson. The
// constituents are retained in flowOrigin so the flow can be unmerged (fallback)
// if no shared slot can be found.
func mergeFlow(constituents []*lesson, k flowKey) *lesson {
	primary := constituents[0].groupID
	subLessons := make([]*internalSubLesson, 0, len(constituents))
	for _, c := range constituents {
		if c.groupID < primary {
			primary = c.groupID
		}
		for _, sl := range c.subLessons {
			subLessons = append(subLessons, &internalSubLesson{
				subGroupNumber: nil,
				groupID:        sl.groupID,
				teacherID:      sl.teacherID,
				roomID:         0,
			})
		}
		// Constituents become plain lessons; they only re-enter scheduling via
		// the fallback path (unmerge).
		c.isFlow = false
	}

	return &lesson{
		groupID:          primary,
		subjectID:        k.subjectID,
		subLessons:       subLessons,
		format:           formatUnited,
		isScheduled:      false,
		lType:            lessonTypeParallel,
		workloadID:       constituents[0].workloadID,
		cycleCommitteeID: constituents[0].cycleCommitteeID,
		isLab:            false,
		flowEligible:     false,
		isFlow:           true,
		flowID:           fmt.Sprintf("flow-%d-%d", k.subjectID, k.teacherID),
		flowOrigin:       constituents,
	}
}

// flowRoomFeasible reports whether at least one room can hold the combined size of
// all groups in a flow lesson. A flow whose groups can never share a room is not
// worth forming and should fall back to separate per-group lessons.
//
// It is permissive when sizes/capacities are unknown: an unknown combined size or
// the presence of an unlimited room (capacity 0) means "no constraint".
func flowRoomFeasible(cfg *settings, l *lesson) bool {
	if cfg == nil || len(cfg.roomCapacities) == 0 {
		return true
	}

	required := 0
	for _, gid := range l.groupIDs() {
		size := cfg.groupSizes[gid]
		if size <= 0 {
			return true // unknown combined size → no hard constraint
		}
		required += size
	}
	if required <= 0 {
		return true
	}

	maxCapacity := 0
	for _, c := range cfg.roomCapacities {
		if c <= 0 {
			return true // an unlimited room fits anything
		}
		if c > maxCapacity {
			maxCapacity = c
		}
	}
	return maxCapacity >= required
}

// partitionFlowLessons splits a lesson list into flow lessons and the rest.
func partitionFlowLessons(lessons []*lesson) (flow []*lesson, regular []*lesson) {
	for _, l := range lessons {
		if l.isFlow {
			flow = append(flow, l)
		} else {
			regular = append(regular, l)
		}
	}
	return flow, regular
}

// sortedKeys returns the map keys in ascending order for deterministic iteration.
func sortedKeys(m map[uint64][]*lesson) []uint64 {
	keys := make([]uint64, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}

// flowDiagnostics accumulates information about flow scheduling for one generation
// attempt: how many flow lessons were formed and which ones had to fall back to
// separate per-group lessons because no shared slot/room was available.
type flowDiagnostics struct {
	flowLessons   int
	fallbackIDs   map[string]bool
	fallbackOrder []string
}

// newFlowDiagnostics creates an empty diagnostics collector.
func newFlowDiagnostics() *flowDiagnostics {
	return &flowDiagnostics{
		flowLessons:   0,
		fallbackIDs:   make(map[string]bool),
		fallbackOrder: nil,
	}
}

// countFlows records how many flow lessons were formed (across both week periods).
func (d *flowDiagnostics) countFlows(lessons ...[]*lesson) {
	if d == nil {
		return
	}
	for _, list := range lessons {
		for _, l := range list {
			if l.isFlow {
				d.flowLessons++
			}
		}
	}
}

// recordFallback notes that a flow lesson fell back to separate per-group lessons.
func (d *flowDiagnostics) recordFallback(l *lesson) {
	if d == nil || l == nil {
		return
	}
	if !d.fallbackIDs[l.flowID] {
		d.fallbackIDs[l.flowID] = true
		d.fallbackOrder = append(d.fallbackOrder, l.flowID)
	}
}

// metadata serializes the diagnostics into a compact JSON string suitable for
// ScheduleData.Metadata. Returns an empty string when no flow lessons were formed.
func (d *flowDiagnostics) metadata() string {
	if d == nil || d.flowLessons == 0 {
		return ""
	}
	payload := struct {
		FlowLessons   int      `json:"flow_lessons"`
		FlowFallbacks int      `json:"flow_fallbacks"`
		FallbackIDs   []string `json:"fallback_flow_ids,omitempty"`
	}{
		FlowLessons:   d.flowLessons,
		FlowFallbacks: len(d.fallbackOrder),
		FallbackIDs:   d.fallbackOrder,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(b)
}
