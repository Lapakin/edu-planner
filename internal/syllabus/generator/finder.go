package generator

import (
	"sort"
	"strings"
	"sync"
)

// finder finds available lesson slots for groups, with caching.
type finder struct {
	mu    sync.RWMutex
	cache map[string][]int // cache key → available lesson numbers
	cfg   *settings
}

// newFinder creates a new finder with the given settings.
func newFinder(cfg *settings) *finder {
	return &finder{
		mu:    sync.RWMutex{},
		cache: make(map[string][]int),
		cfg:   cfg,
	}
}

// findUnscheduledNumbers returns available lesson numbers given the busy slots.
// Results are cached for performance.
func (f *finder) findUnscheduledNumbers(busySlots []int) []int {
	sort.Ints(busySlots)
	key := f.buildCacheKey(busySlots)

	// Check cache (read lock)
	f.mu.RLock()
	if cached, ok := f.cache[key]; ok {
		f.mu.RUnlock()
		return cached
	}
	f.mu.RUnlock()

	// Compute available slots
	available := f.compute(busySlots)

	// Store in cache (write lock)
	f.mu.Lock()
	f.cache[key] = available
	f.mu.Unlock()

	return available
}

// compute calculates available lesson numbers that are not in busy slots
// and are adjacent to existing busy slots (to minimize gaps).
func (f *finder) compute(busySlots []int) []int {
	if len(busySlots) == 0 {
		// If nothing is scheduled, return the first lesson number
		return []int{f.cfg.firstLessonNumber}
	}

	busySet := make(map[int]bool, len(busySlots))
	for _, slot := range busySlots {
		busySet[slot] = true
	}

	available := make([]int, 0)

	// Find slots adjacent to existing busy slots first (minimize gaps)
	for _, busy := range busySlots {
		// Check slot before
		before := busy - 1
		if before >= f.cfg.firstLessonNumber && !busySet[before] {
			available = append(available, before)
		}
		// Check slot after
		after := busy + 1
		if after <= f.cfg.lastLessonNumber && !busySet[after] {
			available = append(available, after)
		}
	}

	// Deduplicate
	seen := make(map[int]bool)
	deduped := make([]int, 0, len(available))
	for _, n := range available {
		if !seen[n] {
			seen[n] = true
			deduped = append(deduped, n)
		}
	}

	sort.Ints(deduped)

	// If no adjacent slots, return all free slots
	if len(deduped) == 0 {
		for i := f.cfg.firstLessonNumber; i <= f.cfg.lastLessonNumber; i++ {
			if !busySet[i] {
				deduped = append(deduped, i)
			}
		}
	}

	return deduped
}

// buildCacheKey creates a string key from sorted busy slots.
func (f *finder) buildCacheKey(slots []int) string {
	parts := make([]string, len(slots))
	for i, s := range slots {
		parts[i] = string(rune('0' + s))
	}
	return strings.Join(parts, ",")
}

// findAllFreeSlots returns all lesson numbers not occupied by any of the given
// groups on the date. Passing several groups (a flow lesson) yields slots free for
// all of them.
func (f *finder) findAllFreeSlots(sched *schedule, d date, groupIDs []uint64) []int {
	busy := sched.getScheduledLessonNumbersForGroups(d, groupIDs)
	busySet := make(map[int]bool, len(busy))
	for _, ln := range busy {
		busySet[ln] = true
	}

	free := make([]int, 0)
	for i := f.cfg.firstLessonNumber; i <= f.cfg.lastLessonNumber; i++ {
		if !busySet[i] {
			free = append(free, i)
		}
	}
	return free
}
