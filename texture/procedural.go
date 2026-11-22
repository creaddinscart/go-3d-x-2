package texture

import (
	"math"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
)

// Solid returns a 1x1 texture filled with one color.
func Solid(c m32.Color) *Texture {
	t := New(1, 1)
	t.Pix[0] = c
	return t
}

// Checker builds a checkerboard of two colors.
func Checker(cols, rows int, a, b m32.Color) *Texture {
	t := New(cols*8, rows*8)
	for y := 0; y < t.H; y++ {
		for x := 0; x < t.W; x++ {
			if ((x/8)+(y/8))%2 == 0 {
				t.Set(x, y, a)
			} else {
				t.Set(x, y, b)
			}
		}
	}
	return t
}

// Gradient builds a horizontal linear gradient between two colors.
func Gradient(w int, a, b m32.Color) *Texture {
	t := New(w, 1)
	for x := 0; x < w; x++ {
		t.Set(x, 0, a.Mix(b, float32(x)/float32(w-1)))
	}
	return t
}

func hash2(x, y int) float32 {
	h := uint32(x*374761393 + y*668265263)
	h = (h ^ (h >> 13)) * 1274126177
	h ^= h >> 16
	return float32(h&0xffffff) / float32(0xffffff)
}

func valueNoise(x, y float32) float32 {
	ix, iy := int(math.Floor(float64(x))), int(math.Floor(float64(y)))
	fx := x - float32(ix)
	fy := y - float32(iy)
	sx := fx * fx * (3 - 2*fx)
	sy := fy * fy * (3 - 2*fy)
	a := hash2(ix, iy)
	b := hash2(ix+1, iy)
	c := hash2(ix, iy+1)
	d := hash2(ix+1, iy+1)
	return m32.Lerp(m32.Lerp(a, b, sx), m32.Lerp(c, d, sx), sy)
}

// FBM returns fractal Brownian motion noise in [0, 1].
func FBM(x, y float32, octaves int) float32 {
	sum, amp, freq := float32(0), float32(0.5), float32(1)
	for i := 0; i < octaves; i++ {
		sum += amp * valueNoise(x*freq, y*freq)
		freq *= 2
		amp *= 0.5
	}
	return m32.Clamp01(sum)
}

// Noise builds a grayscale value-noise texture.
func Noise(w, h int, scale float32, octaves int) *Texture {
	t := New(w, h)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			n := FBM(float32(x)/float32(w)*scale, float32(y)/float32(h)*scale, octaves)
			t.Set(x, y, m32.RGB(n, n, n))
		}
	}
	return t
}

// Brick builds a simple brick-wall texture.
func Brick(w, h int, mortar, brick m32.Color) *Texture {
	t := New(w, h)
	bw, bh, gap := w/4, h/8, maxInt(1, w/64)
	for y := 0; y < h; y++ {
		row := y / bh
		off := (row % 2) * (bw / 2)
		for x := 0; x < w; x++ {
			xx := (x + off) % bw
			if y%bh < gap || xx < gap {
				t.Set(x, y, mortar)
			} else {
				t.Set(x, y, brick.Mul(0.8+0.3*hash2(x/bw, y/bh)))
			}
		}
	}
	return t
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
