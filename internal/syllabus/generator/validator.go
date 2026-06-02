package generator

import (
	"fmt"
	"sort"
)

// validator performs pre-generation validation to ensure scheduling is possible.
type validator struct {
	cfg *settings
}

// newValidator creates a new validator.
func newValidator(cfg *settings) *validator {
	return &validator{cfg: cfg}
}

// validate checks that the generation is possible given the lessons and available dates.
func (v *validator) validate(
	lessons []*lesson,
	dates []date,
	groupIDs []uint64,
) error {
	// Count required slots per group. A flow lesson is attended by every group in
	// its stream, so it counts once toward each of them.
	groupLessonCount := make(map[uint64]int)
	for _, l := range lessons {
		for _, gid := range l.groupIDs() {
			groupLessonCount[gid]++
		}
	}

	// Calculate available slots per group (slots per day * number of days)
	daysCount := len(dates)
	maxSlotsPerGroup := daysCount * v.cfg.maxLessonsPerDay()

	for _, gid := range groupIDs {
		required := groupLessonCount[gid]
		if required > maxSlotsPerGroup {
			return fmt.Errorf(
				"%w: group %d requires %d lesson slots but only %d are available (%d days × %d slots/day)",
				ErrValidationFailed, gid, required, maxSlotsPerGroup, daysCount, v.cfg.maxLessonsPerDay(),
			)
		}
	}

	// Check weekly limit
	maxGroupPerWeek := v.cfg.maxGroupLessonsPerWeek()
	for _, gid := range groupIDs {
		required := groupLessonCount[gid]
		if required > maxGroupPerWeek {
			return fmt.Errorf(
				"%w: group %d requires %d lessons/week but max is %d",
				ErrValidationFailed, gid, required, maxGroupPerWeek,
			)
		}
	}

	// Capacity feasibility: every group that has lessons and a known size must fit
	// in at least one room. A room with capacity 0 is treated as unlimited, so the
	// check only fails when every room has a positive capacity smaller than the group.
	if len(v.cfg.roomCapacities) > 0 && len(v.cfg.groupSizes) > 0 {
		maxRoomCapacity := 0
		hasUnlimitedRoom := false
		for _, c := range v.cfg.roomCapacities {
			if c <= 0 {
				hasUnlimitedRoom = true
				break
			}
			if c > maxRoomCapacity {
				maxRoomCapacity = c
			}
		}

		if !hasUnlimitedRoom {
			for _, gid := range groupIDs {
				if groupLessonCount[gid] == 0 {
					continue
				}
				size := v.cfg.groupSizes[gid]
				if size > 0 && size > maxRoomCapacity {
					return fmt.Errorf(
						"%w: group %d has %d students but the largest room holds only %d",
						ErrValidationFailed, gid, size, maxRoomCapacity,
					)
				}
			}
		}
	}

	// Check minimum lessons per day: if a group has any lessons but fewer than min,
	// it's impossible to satisfy the constraint on any day.
	if v.cfg.minGroupLessonsPerDay > 0 {
		for _, gid := range groupIDs {
			required := groupLessonCount[gid]
			if required > 0 && required < v.cfg.minGroupLessonsPerDay {
				return fmt.Errorf(
					"%w: group %d requires %d lesson(s)/week but min per day is %d (cannot satisfy daily minimum)",
					ErrValidationFailed, gid, required, v.cfg.minGroupLessonsPerDay,
				)
			}
		}
	}

	return nil
}

// validatePlacement checks post-generation constraints: no gaps and minimum lessons per day.
// Returns ErrUnableToSchedule if the placed schedule violates hard constraints.
func validatePlacement(sched *schedule, dates []date, groupIDs []uint64, cfg *settings) error {
	for _, d := range dates {
		for _, gid := range groupIDs {
			numbers := sched.getScheduledLessonNumbers(d, gid)
			if len(numbers) == 0 {
				continue
			}

			if cfg.minGroupLessonsPerDay > 0 && len(numbers) < cfg.minGroupLessonsPerDay {
				return fmt.Errorf(
					"%w: group %d has %d lesson(s) on a day (minimum %d)",
					ErrUnableToSchedule, gid, len(numbers), cfg.minGroupLessonsPerDay,
				)
			}

			if cfg.noGapsRequired {
				sort.Ints(numbers)
				for i := 1; i < len(numbers); i++ {
					if numbers[i]-numbers[i-1] > 1 {
						return fmt.Errorf(
							"%w: group %d has a gap between lesson %d and %d",
							ErrUnableToSchedule, gid, numbers[i-1], numbers[i],
						)
					}
				}
			}
		}
	}
	return nil
}
