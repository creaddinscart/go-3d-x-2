// Package texture provides an RGBA texture with bilinear sampling and wrap
// modes, conversion to/from standard-library images, and procedural texture
// generators for 3d-x-2.
//
// Built by CREADDINSCART STUDIO — https://www.shit.pub/studio/
package texture

import (
	"image"
	"image/color"
	"math"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
)

// WrapMode controls texture coordinate addressing outside [0, 1].
type WrapMode int

const (
	Repeat WrapMode = iota
	Clamp
	Mirror
)

// Texture is an RGBA float32 texture sampled with bilinear filtering.
type Texture struct {
	W, H    int
	Pix     []m32.Color
	Wrap    WrapMode
	Nearest bool
}

// New allocates a texture of the given size.
func New(w, h int) *Texture {
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	return &Texture{W: w, H: h, Pix: make([]m32.Color, w*h), Wrap: Repeat}
}

// Set writes a single texel.
func (t *Texture) Set(x, y int, c m32.Color) {
	if x >= 0 && x < t.W && y >= 0 && y < t.H {
		t.Pix[y*t.W+x] = c
	}
}

// At reads a single texel (clamped).
func (t *Texture) At(x, y int) m32.Color {
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	if x >= t.W {
		x = t.W - 1
	}
	if y >= t.H {
		y = t.H - 1
	}
	return t.Pix[y*t.W+x]
}

func (t *Texture) wrapCoord(i, size int) int {
	switch t.Wrap {
	case Clamp:
		return int(m32.Clamp(float32(i), 0, float32(size-1)))
	case Mirror:
		m := i % (2 * size)
		if m < 0 {
			m += 2 * size
		}
		if m >= size {
			m = 2*size - 1 - m
		}
		return m
	default:
		m := i % size
		if m < 0 {
			m += size
		}
		return m
	}
}

// Sample evaluates the texture at continuous UV coordinates.
func (t *Texture) Sample(u, v float32) m32.Color {
	if t.W == 0 || t.H == 0 {
		return m32.White
	}
	fx := u*float32(t.W) - 0.5
	fy := v*float32(t.H) - 0.5
	x0 := int(math.Floor(float64(fx)))
	y0 := int(math.Floor(float64(fy)))
	if t.Nearest {
		return t.At(t.wrapCoord(x0, t.W), t.wrapCoord(y0, t.H))
	}
	tx := float32(fx) - float32(x0)
	ty := float32(fy) - float32(y0)
	c00 := t.At(t.wrapCoord(x0, t.W), t.wrapCoord(y0, t.H))
	c10 := t.At(t.wrapCoord(x0+1, t.W), t.wrapCoord(y0, t.H))
	c01 := t.At(t.wrapCoord(x0, t.W), t.wrapCoord(y0+1, t.H))
	c11 := t.At(t.wrapCoord(x0+1, t.W), t.wrapCoord(y0+1, t.H))
	return c00.Mix(c10, tx).Mix(c01.Mix(c11, tx), ty)
}

// Fill paints the whole texture with a color.
func (t *Texture) Fill(c m32.Color) {
	for i := range t.Pix {
		t.Pix[i] = c
	}
}

// FromImage builds a texture from any standard-library image.
func FromImage(img image.Image) *Texture {
	b := img.Bounds()
	t := New(b.Dx(), b.Dy())
	for y := 0; y < b.Dy(); y++ {
		for x := 0; x < b.Dx(); x++ {
			c := color.NRGBAModel.Convert(img.At(b.Min.X+x, b.Min.Y+y)).(color.NRGBA)
			t.Pix[y*t.W+x] = m32.FromNRGBA(c)
		}
	}
	return t
}

// ToImage converts the texture to an image.NRGBA.
func (t *Texture) ToImage() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, t.W, t.H))
	for y := 0; y < t.H; y++ {
		for x := 0; x < t.W; x++ {
			img.SetNRGBA(x, y, t.Pix[y*t.W+x].ToColorNRGBA())
		}
	}
	return img
}
