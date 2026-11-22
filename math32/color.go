package math32

import (
	"image/color"
	"math"
)

// Color is an RGBA float32 color with channels in [0, 1].
type Color struct{ R, G, B, A float32 }

// RGBA returns the color as 16-bit image/color components.
func (c Color) RGBA() (r, g, b, a uint32) {
	to16 := func(v float32) uint32 { return uint32(Clamp01(v)*65535 + 0.5) }
	return to16(c.R), to16(c.G), to16(c.B), to16(c.A)
}

// ToColorNRGBA converts to a standard library 8-bit NRGBA color.
func (c Color) ToColorNRGBA() color.NRGBA {
	to8 := func(v float32) uint8 { return uint8(Clamp01(v)*255 + 0.5) }
	return color.NRGBA{R: to8(c.R), G: to8(c.G), B: to8(c.B), A: to8(c.A)}
}

// RGB builds an opaque color from float components.
func RGB(r, g, b float32) Color { return Color{r, g, b, 1} }

// NewRGBA builds a color with alpha.
func NewRGBA(r, g, b, a float32) Color { return Color{r, g, b, a} }

// FromNRGBA converts a standard library NRGBA color.
func FromNRGBA(c color.NRGBA) Color {
	return Color{float32(c.R) / 255, float32(c.G) / 255, float32(c.B) / 255, float32(c.A) / 255}
}

// Hex parses a #RGB, #RGBA, #RRGGBB or #RRGGBBAA literal.
func Hex(s string) Color {
	var r, g, b, a int = 255, 255, 255, 255
	if len(s) > 0 && s[0] == '#' {
		s = s[1:]
	}
	switch len(s) {
	case 3:
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2]})
	case 4:
		s = string([]byte{s[0], s[0], s[1], s[1], s[2], s[2], s[3], s[3]})
	}
	n := 0
	for i := 0; i < len(s) && i < 8; i++ {
		n <<= 4
		switch {
		case s[i] >= '0' && s[i] <= '9':
			n |= int(s[i] - '0')
		case s[i] >= 'a' && s[i] <= 'f':
			n |= int(s[i]-'a') + 10
		case s[i] >= 'A' && s[i] <= 'F':
			n |= int(s[i]-'A') + 10
		}
	}
	switch len(s) {
	case 8:
		a = n & 0xff
		n >>= 8
	}
	b = n & 0xff
	n >>= 8
	g = n & 0xff
	n >>= 8
	r = n & 0xff
	return Color{float32(r) / 255, float32(g) / 255, float32(b) / 255, float32(a) / 255}
}

// Mul scales every channel of the color by k.
func (c Color) Mul(k float32) Color {
	return Color{c.R * k, c.G * k, c.B * k, c.A}
}

// Modulate multiplies two colors component-wise.
func (c Color) Modulate(o Color) Color {
	return Color{c.R * o.R, c.G * o.G, c.B * o.B, c.A * o.A}
}

// AddColor adds two colors, clamping each channel to [0, 1].
func (c Color) AddColor(o Color) Color {
	return Color{Clamp01(c.R + o.R), Clamp01(c.G + o.G), Clamp01(c.B + o.B), Clamp01(c.A + o.A)}
}

// Mix linearly blends c toward o by t.
func (c Color) Mix(o Color, t float32) Color {
	return Color{Lerp(c.R, o.R, t), Lerp(c.G, o.G, t), Lerp(c.B, o.B, t), Lerp(c.A, o.A, t)}
}

// HSV builds a color from hue [0,360), saturation and value in [0,1].
func HSV(h, s, v float32) Color {
	h = float32(math.Mod(float64(h), 360))
	hi := int32(h / 60)
	f := h/60 - float32(hi)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	switch hi {
	case 0:
		return RGB(v, t, p)
	case 1:
		return RGB(q, v, p)
	case 2:
		return RGB(p, v, t)
	case 3:
		return RGB(p, q, v)
	case 4:
		return RGB(t, p, v)
	default:
		return RGB(v, p, q)
	}
}

// Common palette constants.
var (
	Black       = RGB(0, 0, 0)
	White       = RGB(1, 1, 1)
	Red         = RGB(1, 0, 0)
	Green       = RGB(0, 1, 0)
	Blue        = RGB(0, 0, 1)
	Yellow      = RGB(1, 1, 0)
	Cyan        = RGB(0, 1, 1)
	Magenta     = RGB(1, 0, 1)
	Transparent = Color{0, 0, 0, 0}
)
