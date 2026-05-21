package generator

// coordinator holds all sub-components used during a single generation attempt.
type coordinator struct {
	finder       *finder
	fixer        *fixer
	schedule     *schedule
	availability *availability
	roomSelector *roomSelector
	patternCtrl  *patternController
	limiter      *limiter
	cfg          *settings
}

// newCoordinator creates a new coordinator with all sub-components initialized.
func newCoordinator(cfg *settings, roomIDs []uint64) *coordinator {
	sched := newSchedule()
	avail := newAvailability()
	f := newFinder(cfg)
	patternCtrl := newPatternController()
	lim := newLimiter(cfg)
	rs := newRoomSelector(roomIDs, cfg.cycleCommitteeLabRooms)

	coord := &coordinator{
		finder:       f,
		fixer:        nil,
		schedule:     sched,
		availability: avail,
		roomSelector: rs,
		patternCtrl:  patternCtrl,
		limiter:      lim,
		cfg:          cfg,
	}

	// fixer needs coordinator reference
	coord.fixer = newFixer(cfg, coord)

	return coord
}
