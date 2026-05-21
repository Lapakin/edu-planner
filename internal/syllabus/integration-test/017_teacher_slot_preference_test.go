package integration_test

import (
	"net/http"
	"testing"

	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/utils"

	ta "github.com/Lapakin/edu-planner/internal/testassets"
)

func TestCreateTeacherSlotPreferences(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teacher-slot-preferences", http.MethodPost, adminToken).
		NewOKRequestWithBody(ta.TeacherSlotPreferencesArray).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestGetTeacherSlotPreferenceByID(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teacher-slot-preferences/{teacherSlotPreferenceId}", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, map[string]any{"teacherSlotPreferenceId": ta.TeacherSlotPreference1.ID}, ta.TeacherSlotPreference1).
		NewRequestWithoutBody("NotFound", nil, map[string]any{"teacherSlotPreferenceId": uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestFetchTeacherSlotPreferences(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teacher-slot-preferences", http.MethodGet, adminToken).
		NewOKRequestWithoutBody(nil, nil, ta.TeacherSlotPreferencesArray).
		NewOKRequestWithoutBody(map[string]any{domain.IDsParam: []uint64{ta.TeacherSlotPreference1.ID}}, nil, domain.TeacherSlotPreferences{ta.TeacherSlotPreference1}).
		NewOKRequestWithoutBody(map[string]any{domain.AcademicYearIDParam: ta.AcademicYear1.ID}, nil, ta.TeacherSlotPreferencesArray).
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestUpdateTeacherSlotPreferences(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teacher-slot-preferences", http.MethodPut, adminToken).
		NewRequestWithDataManipulation("OK", ta.TeacherSlotPreferencesArray,
			func(input any) any {
				var modified domain.TeacherSlotPreferences
				utils.Copy(input, &modified)
				modified[0].LessonNumber = 3
				modified[0].SlotType = domain.SlotTypeForbidden
				return modified
			}, http.StatusOK, ta.TeacherSlotPreferencesArray,
		).
		NewRequestWithBody("NotFound", domain.TeacherSlotPreferences{&domain.TeacherSlotPreference{
			ID: 0, TeacherID: ta.Teacher1.ID, Weekday: domain.WeekdayMonday, LessonNumber: 1, SlotType: domain.SlotTypePreferred,
		}}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}

func TestDeleteTeacherSlotPreferences(t *testing.T) {
	ta.NewHTTPCasesRunner("/api/v1/teacher-slot-preferences/delete", http.MethodPost, adminToken).
		NewRequestWithBody("OK", []uint64{ta.TeacherSlotPreference2.ID}, http.StatusOK, nil).
		NewRequestWithBody("NotFound", []uint64{uint64(0)}, http.StatusNotFound, ta.ExpectedResponseNotFound).
		NewBadJWTRequest().
		NewUnmarshalErrorRequest().
		NewNilBodyRequest().
		Run(t, ts.URL, ta.DefaultIgnoredFields)
}
