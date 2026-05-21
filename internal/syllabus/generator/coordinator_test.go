package generator

import "testing"

func TestNewCoordinator_AllFieldsNonNil(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4)
	cfg := newSettings(ts, bells, nil, nil, nil)
	roomIDs := []uint64{101, 102}

	coord := newCoordinator(cfg, roomIDs)

	if coord == nil {
		t.Fatal("expected non-nil coordinator")
	}
	if coord.finder == nil {
		t.Error("expected non-nil finder")
	}
	if coord.fixer == nil {
		t.Error("expected non-nil fixer")
	}
	if coord.schedule == nil {
		t.Error("expected non-nil schedule")
	}
	if coord.availability == nil {
		t.Error("expected non-nil availability")
	}
	if coord.roomSelector == nil {
		t.Error("expected non-nil roomSelector")
	}
	if coord.patternCtrl == nil {
		t.Error("expected non-nil patternCtrl")
	}
	if coord.limiter == nil {
		t.Error("expected non-nil limiter")
	}
	if coord.cfg == nil {
		t.Error("expected non-nil cfg")
	}
}

func TestNewCoordinator_EmptyRoomIDs(t *testing.T) {
	ts := makeTemplateSetting(2, 8, 40)
	bells := makeBellSchedules(4)
	cfg := newSettings(ts, bells, nil, nil, nil)

	coord := newCoordinator(cfg, []uint64{})
	if coord == nil {
		t.Fatal("expected non-nil coordinator with empty room IDs")
	}
	if coord.roomSelector == nil {
		t.Error("expected non-nil roomSelector even with empty room IDs")
	}
}
