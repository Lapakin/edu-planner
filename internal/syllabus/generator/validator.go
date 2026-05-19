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
	// Count required slots per group
	groupLessonCount := make(map[uint64]int)
	for _, l := range lessons {
		groupLessonCount[l.groupID]++
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
