package domain

// AcademicYearToUser is the bridge linking an academic year to a user.
type AcademicYearToUser struct {
	AcademicYearID uint64 `json:"academic_year_id" db:"academic_year_id"`
	UserID         uint64 `json:"user_id" db:"user_id"`
}

// AttachUsersReq is the payload for attaching users to an academic year.
type AttachUsersReq struct {
	AcademicYearID uint64   `json:"academic_year_id"`
	UserIDs        []uint64 `json:"user_ids"`
}
