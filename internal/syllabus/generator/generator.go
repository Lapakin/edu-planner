package generator

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/Lapakin/edu-planner/internal/domain"
)

// Generator orchestrates multi-goroutine parallel schedule generation.
type Generator struct {
	cfg       domain.GenerationConfig
	settings  *settings
	workloads []*Workload
	roomIDs   []uint64
	groupIDs  []uint64
	startDate date
	endDate   date
}

// New creates a new Generator.
func New(
	genCfg domain.GenerationConfig,
	templateSetting *domain.ScheduleTemplateSetting,
	bellSchedules domain.BellSchedules,
	restriction *domain.ScheduleRestriction,
	preferences domain.TeacherSlotPreferences,
	labRooms domain.CycleCommitteeLabRooms,
	workloads []*Workload,
	roomIDs []uint64,
	groupIDs []uint64,
	roomCapacities map[uint64]int,
	groupSizes map[uint64]int,
	startDate, endDate time.Time,
) *Generator {
	s := newSettings(templateSetting, bellSchedules, restriction, preferences, labRooms)
	s.roomCapacities = roomCapacities
	s.groupSizes = groupSizes

	return &Generator{
		cfg:       genCfg,
		settings:  s,
		workloads: workloads,
		roomIDs:   roomIDs,
		groupIDs:  groupIDs,
		startDate: newDate(startDate),
		endDate:   newDate(endDate),
	}
}

// Generate launches N parallel goroutines, each running an independent generation
// attempt. Returns the first valid schedule or an error.
func (g *Generator) Generate(ctx context.Context) (*domain.ScheduleData, error) {
	ctx, cancel := context.WithTimeout(ctx, g.cfg.Timeout)
	defer cancel()

	numGoroutines := max(g.cfg.NumberOfGoroutines, 1)
	maxAttempts := max(g.cfg.MaxAttempts, numGoroutines)
	attemptsPerGoroutine := max(maxAttempts/numGoroutines, 1)

	type result struct {
		data *domain.ScheduleData
		err  error
	}

	results := make(chan result, numGoroutines)
	var wg sync.WaitGroup

	for i := range numGoroutines {
		wg.Add(1)
		go func(goroutineIdx int) {
			defer wg.Done()

			var lastErr error
			for attempt := range attemptsPerGoroutine {
				select {
				case <-ctx.Done():
					if lastErr != nil {
						select {
						case results <- result{data: nil, err: lastErr}:
						default:
						}
					}
					return
				default:
				}

				seed := time.Now().UnixNano() + int64(goroutineIdx)*1000 + int64(attempt)*997 + rand.Int63()
				gen := newGeneration(
					g.settings,
					g.workloads,
					g.roomIDs,
					g.groupIDs,
					g.startDate,
					g.endDate,
					seed,
				)

				data, err := gen.exec()
				if err == nil && data != nil {
					select {
					case results <- result{data: data, err: nil}:
					case <-ctx.Done():
					}
					return
				}
				lastErr = err
			}

			// All attempts for this goroutine exhausted
			select {
			case results <- result{data: nil, err: lastErr}:
			case <-ctx.Done():
			}
		}(i)
	}

	// Close results channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	var lastErr error

	for {
		select {
		case <-ctx.Done():
			if lastErr != nil {
				return nil, lastErr
			}
			return nil, ErrTimeout

		case res, ok := <-results:
			if !ok {
				// All goroutines finished without success
				if lastErr != nil {
					return nil, lastErr
				}
				return nil, ErrUnableToSchedule
			}

			if res.err == nil && res.data != nil {
				// Success! Cancel other goroutines and return
				cancel()
				return res.data, nil
			}

			if res.err != nil {
				lastErr = res.err
			}
		}
	}
}

// BuildWorkloads is a public helper to convert domain models to internal Workloads.
func BuildWorkloads(
	distributions domain.WorkloadDistributions,
	assignments domain.WorkloadAssignments,
	studyPlans domain.StudyPlans,
	groups domain.Groups,
	disciplines domain.Disciplines,
) []*Workload {
	return buildWorkloadsFromDomain(distributions, assignments, studyPlans, groups, disciplines)
}
