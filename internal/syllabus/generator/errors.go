package generator

import "errors"

var (
	// ErrTimeout is returned when generation exceeds the configured timeout.
	ErrTimeout = errors.New("schedule generation timed out")

	// ErrWorkloadDistribution is returned when workloads cannot be distributed fairly.
	ErrWorkloadDistribution = errors.New("unable to distribute workloads across weeks")

	// ErrUnableToSchedule is returned when the recursive scheduler cannot place all lessons.
	ErrUnableToSchedule = errors.New("unable to schedule all lessons")

	// ErrUnexpectedBehavior is returned when an unexpected state is reached.
	ErrUnexpectedBehavior = errors.New("unexpected behavior during generation")

	// ErrGroupBusy is returned when a group has a scheduling conflict.
	ErrGroupBusy = errors.New("group has scheduling conflict")

	// ErrTeacherBusy is returned when a teacher has a scheduling conflict.
	ErrTeacherBusy = errors.New("teacher has scheduling conflict")

	// ErrNoSettings is returned when schedule template settings are not found.
	ErrNoSettings = errors.New("schedule template settings not found")

	// ErrNoBellSchedule is returned when bell schedules are not configured.
	ErrNoBellSchedule = errors.New("bell schedule not configured")

	// ErrNoGroupSemesters is returned when no group semesters are found for the given semester.
	ErrNoGroupSemesters = errors.New("no group semesters found for the given semester")

	// ErrNoWorkloads is returned when no workloads are found for the given groups.
	ErrNoWorkloads = errors.New("no workloads found for the given groups")

	// ErrValidationFailed is returned when pre-generation validation fails.
	ErrValidationFailed = errors.New("pre-generation validation failed")
)
