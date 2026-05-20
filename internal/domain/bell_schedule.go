package domain

type BellSchedule struct {
	ID             uint64 `json:"id" db:"id"`
	AcademicYearID uint64 `json:"academic_year_id" db:"academic_year_id"`
	LessonNumber   int    `json:"lesson_number" db:"lesson_number"`
	StartTime      string `json:"start_time" db:"start_time"`
	EndTime        string `json:"end_time" db:"end_time"`
}

type BellSchedules []*BellSchedule
