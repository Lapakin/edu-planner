package domain

import "time"

type GenerateScheduleRequest struct {
	SemesterID uint64   `json:"semester_id"`
	GroupIDs   []uint64 `json:"group_ids,omitempty"`
}

type SaveScheduleTemplateRequest struct {
	SemesterID uint64        `json:"semester_id"`
	Name       *string       `json:"name,omitempty"`
	Data       *ScheduleData `json:"data"`
}

type GenerationConfig struct {
	Timeout            time.Duration
	NumberOfGoroutines int
}

func DefaultGenerationConfig() GenerationConfig {
	return GenerationConfig{
		Timeout:            120 * time.Second,
		NumberOfGoroutines: 4,
	}
}
