package core

import (
    "time"
)

// TickManager handles the tick channels for the service
type TickManager struct {
    FastTick    <-chan time.Time
    SlowTick    <-chan time.Time
    CurrentTick <-chan time.Time
}

// NewTickManager creates a new tick manager with default intervals
func NewTickManager() *TickManager {
    fastTick := time.Tick(2 * time.Second)
    slowTick := time.Tick(5 * time.Second)

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
