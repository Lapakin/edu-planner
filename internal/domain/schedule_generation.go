package domain

import "time"

type GenerateScheduleRequest struct {
	SemesterID uint64 `json:"semester_id"`
}

type SaveScheduleTemplateRequest struct {
	SemesterID uint64        `json:"semester_id"`
	Name       *string       `json:"name,omitempty"`
	Data       *ScheduleData `json:"data"`
}

type GenerationConfig struct {
	Timeout            time.Duration
	NumberOfGoroutines int
	MaxAttempts        int
}

func DefaultGenerationConfig() GenerationConfig {
	return GenerationConfig{
		Timeout:            120 * time.Second,
		NumberOfGoroutines: 4,
		MaxAttempts:        100,
	}
}
