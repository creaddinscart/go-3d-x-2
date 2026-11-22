// Package anim provides the animation layer of 3d-x-2: an easing library,
// scalar tweens, Vec3 keyframe tracks, clips, a scene-driving Animator, a
// fixed-timestep GameLoop and a ParticleSystem that emits billboard points.
//
// Built by CREADDINSCART STUDIO — https://www.shit.pub/studio/
package anim

import "math"

// EaseFunc maps a normalized time t in [0,1] to an eased value.
type EaseFunc func(t float32) float32

// Built-in easing functions.
var (
	EaseLinear    EaseFunc = func(t float32) float32 { return t }
	EaseInQuad    EaseFunc = func(t float32) float32 { return t * t }
	EaseOutQuad   EaseFunc = func(t float32) float32 { return t * (2 - t) }
	EaseInOutQuad EaseFunc = func(t float32) float32 {
		if t < 0.5 {
			return 2 * t * t
		}
		return 1 - float32(math.Pow(float64(-2*t+2), 2))/2
	}
	EaseInCubic    EaseFunc = func(t float32) float32 { return t * t * t }
	EaseOutCubic   EaseFunc = func(t float32) float32 { f := t - 1; return f*f*f + 1 }
	EaseInOutCubic EaseFunc = func(t float32) float32 {
		if t < 0.5 {
			return 4 * t * t * t
		}
		f := -2*t + 2
		return 1 - float32(math.Pow(float64(f), 3))/2
	}
	EaseInSine     EaseFunc = func(t float32) float32 { return 1 - float32(math.Cos(float64(t)*math.Pi/2)) }
	EaseOutSine    EaseFunc = func(t float32) float32 { return float32(math.Sin(float64(t) * math.Pi / 2)) }
	EaseOutElastic EaseFunc = func(t float32) float32 {
		if t == 0 || t == 1 {
			return t
		}
		c4 := 2 * float32(math.Pi) / 3
		return float32(math.Pow(2, float64(-10*t)))*float32(math.Sin(float64((t*10-0.75)*c4))) + 1
	}
	EaseOutBounce EaseFunc = func(t float32) float32 {
		n1, d1 := float32(7.5625), float32(2.75)
		switch {
		case t < 1/d1:
			return n1 * t * t
		case t < 2/d1:
			t -= 1.5 / d1
			return n1*t*t + 0.75
		case t < 2.5/d1:
			t -= 2.25 / d1
			return n1*t*t + 0.9375
		default:
			t -= 2.625 / d1
			return n1*t*t + 0.984375
		}
	}
)
