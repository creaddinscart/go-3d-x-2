// Package math32 is the linear-algebra core of 3d-x-2.
//
// It is a leaf package with no internal dependencies: column vectors, a
// column-major 4x4 matrix, quaternions, colors and scalar helpers. Everything
// else in the library builds on these types.
//
// Derived from CREADDINSCART STUDIO's 3d-x project.
// Built by CREADDINSCART STUDIO — https://www.shit.pub/studio/
package math32

import "math"

// Clamp restricts v to the inclusive range [lo, hi].
func Clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// Clamp01 restricts v to [0, 1].
func Clamp01(v float32) float32 { return Clamp(v, 0, 1) }

// Lerp linearly interpolates between a and b by t.
func Lerp(a, b, t float32) float32 { return a + (b-a)*t }

// InvLerp returns the normalized position of v within [a, b].
func InvLerp(a, b, v float32) float32 {
	if a == b {
		return 0
	}
	return (v - a) / (b - a)
}

// Remap maps v from [inA, inB] into [outA, outB].
func Remap(v, inA, inB, outA, outB float32) float32 {
	return Lerp(outA, outB, InvLerp(inA, inB, v))
}

// SmoothStep is the Hermite smooth interpolation of t.
func SmoothStep(t float32) float32 { t = Clamp01(t); return t * t * (3 - 2*t) }

// Deg2Rad converts degrees to radians.
func Deg2Rad(d float32) float32 { return d * float32(math.Pi) / 180 }

// Rad2Deg converts radians to degrees.
func Rad2Deg(r float32) float32 { return r * 180 / float32(math.Pi) }

// Min returns the smaller of a and b.
func Min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of a and b.
func Max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of a.
func Abs(a float32) float32 {
	if a < 0 {
		return -a
	}
	return a
}

// Min3 returns the smallest of three floats.
func Min3(a, b, c float32) float32 { return Min(Min(a, b), c) }

// Max3 returns the largest of three floats.
func Max3(a, b, c float32) float32 { return Max(Max(a, b), c) }
