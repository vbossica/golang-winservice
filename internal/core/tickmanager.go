package core

import (
	"time"
)

type TickManager struct {
	FastTick    <-chan time.Time
	SlowTick    <-chan time.Time
	CurrentTick <-chan time.Time
}

func NewTickManager(fastTickValue int, slowTickValue int) *TickManager {
	fastTick := time.Tick(time.Duration(fastTickValue) * time.Second)
	slowTick := time.Tick(time.Duration(slowTickValue) * time.Second)

	return &TickManager{
		FastTick:    fastTick,
		SlowTick:    slowTick,
		CurrentTick: fastTick, // Default to fast tick
	}
}

// UseFastTick switches the current tick to fast tick
func (tm *TickManager) UseFastTick() {
	tm.CurrentTick = tm.FastTick
}

// UseSlowTick switches the current tick to slow tick
func (tm *TickManager) UseSlowTick() {
	tm.CurrentTick = tm.SlowTick
}
