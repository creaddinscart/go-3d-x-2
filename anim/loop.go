package anim

import "time"

// GameLoop runs a fixed-timestep update with a render callback ("fix your
// timestep").
type GameLoop struct {
	TargetFPS float32
	MaxFrames int
	FixedStep time.Duration
}

// NewGameLoop builds a loop at the given frame rate.
func NewGameLoop(targetFPS float32) *GameLoop {
	step := time.Duration(float32(time.Second) / targetFPS)
	return &GameLoop{TargetFPS: targetFPS, FixedStep: step}
}

// Run drives update(dtSeconds) and render() until maxFrames is reached or the
// render callback returns false. update may be nil.
func (g *GameLoop) Run(update func(dt float32), render func() bool) {
	tick := time.NewTicker(g.FixedStep)
	defer tick.Stop()
	frames := 0
	for range tick.C {
		if update != nil {
			update(float32(g.FixedStep) / float32(time.Second))
		}
		if render != nil && !render() {
			return
		}
		frames++
		if g.MaxFrames > 0 && frames >= g.MaxFrames {
			return
		}
	}
}
