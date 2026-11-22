// Package render implements 3d-x-2's software rasterizer and framebuffer. It
// turns a scene into pixels using an incremental (edge-function) triangle
// rasterizer with a z-buffer, perspective-correct interpolation, Blinn-Phong
// lighting, fog and optional supersampling anti-aliasing. Output can be an
// image (PNG / animated GIF) or ANSI truecolor art for a terminal.
//
// Built by CREADDINSCART STUDIO — https://www.shit.pub/studio/
package render

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"io"
	"os"
	"strconv"
	"strings"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
)

// Framebuffer is a color + depth buffer that a Renderer draws into.
type Framebuffer struct {
	W, H  int
	Pix   []m32.Color
	Depth []float32
}

// NewFramebuffer allocates a cleared framebuffer.
func NewFramebuffer(w, h int) *Framebuffer {
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	fb := &Framebuffer{W: w, H: h, Pix: make([]m32.Color, w*h), Depth: make([]float32, w*h)}
	fb.Clear(m32.Black, 1)
	return fb
}

// Clear fills the color and depth buffers.
func (fb *Framebuffer) Clear(c m32.Color, depth float32) {
	for i := range fb.Pix {
		fb.Pix[i] = c
		fb.Depth[i] = depth
	}
}

// Fill paints the color buffer (depth preserved).
func (fb *Framebuffer) Fill(c m32.Color) {
	for i := range fb.Pix {
		fb.Pix[i] = c
	}
}

// SetPixel writes a color unconditionally.
func (fb *Framebuffer) SetPixel(x, y int, c m32.Color) {
	if x >= 0 && x < fb.W && y >= 0 && y < fb.H {
		fb.Pix[y*fb.W+x] = c
	}
}

// TryDepth writes the pixel if z passes the depth test (smaller is nearer).
func (fb *Framebuffer) TryDepth(x, y int, z float32, c m32.Color) bool {
	if x < 0 || x >= fb.W || y < 0 || y >= fb.H {
		return false
	}
	idx := y*fb.W + x
	if z < fb.Depth[idx] {
		fb.Depth[idx] = z
		fb.Pix[idx] = c
		return true
	}
	return false
}

// Blend writes a pixel honoring its alpha (over the existing color).
func (fb *Framebuffer) Blend(x, y int, c m32.Color) {
	if x < 0 || x >= fb.W || y < 0 || y >= fb.H {
		return
	}
	idx := y*fb.W + x
	dst := fb.Pix[idx]
	a := c.A
	fb.Pix[idx] = m32.NewRGBA(
		m32.Lerp(dst.R, c.R, a), m32.Lerp(dst.G, c.G, a),
		m32.Lerp(dst.B, c.B, a), m32.Clamp01(dst.A+a),
	)
}

// Add writes a pixel additively (used for glowing particles).
func (fb *Framebuffer) Add(x, y int, c m32.Color) {
	if x < 0 || x >= fb.W || y < 0 || y >= fb.H {
		return
	}
	idx := y*fb.W + x
	dst := fb.Pix[idx]
	fb.Pix[idx] = m32.NewRGBA(m32.Clamp01(dst.R+c.R), m32.Clamp01(dst.G+c.G), m32.Clamp01(dst.B+c.B), 1)
}

// Pixel reads a color.
func (fb *Framebuffer) Pixel(x, y int) m32.Color { return fb.Pix[y*fb.W+x] }

// ToImage converts the framebuffer to an image.NRGBA.
func (fb *Framebuffer) ToImage() *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, fb.W, fb.H))
	for y := 0; y < fb.H; y++ {
		for x := 0; x < fb.W; x++ {
			img.SetNRGBA(x, y, fb.Pix[y*fb.W+x].ToColorNRGBA())
		}
	}
	return img
}

// SavePNG writes the framebuffer to a PNG file.
func (fb *Framebuffer) SavePNG(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, fb.ToImage())
}

// WritePNG encodes the framebuffer as PNG to any writer.
func (fb *Framebuffer) WritePNG(w io.Writer) error { return png.Encode(w, fb.ToImage()) }

func quant8(v float32) int {
	n := int(v*255 + 0.5)
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return n
}

// ToANSI renders the framebuffer as 24-bit truecolor half-block art: each
// terminal cell packs two vertical pixels (background = top, foreground =
// bottom) with the "▀" glyph, producing H/2 rows of W columns.
func (fb *Framebuffer) ToANSI() string {
	var b strings.Builder
	b.Grow(fb.W * (fb.H/2 + 1) * 22)
	for y := 0; y < fb.H; y += 2 {
		for x := 0; x < fb.W; x++ {
			top := fb.Pix[y*fb.W+x]
			bot := top
			if y+1 < fb.H {
				bot = fb.Pix[(y+1)*fb.W+x]
			}
			b.WriteString("\x1b[48;2;")
			b.WriteString(strconv.Itoa(quant8(top.R)))
			b.WriteByte(';')
			b.WriteString(strconv.Itoa(quant8(top.G)))
			b.WriteByte(';')
			b.WriteString(strconv.Itoa(quant8(top.B)))
			b.WriteString("m\x1b[38;2;")
			b.WriteString(strconv.Itoa(quant8(bot.R)))
			b.WriteByte(';')
			b.WriteString(strconv.Itoa(quant8(bot.G)))
			b.WriteByte(';')
			b.WriteString(strconv.Itoa(quant8(bot.B)))
			b.WriteString("m▀")
		}
		b.WriteString("\x1b[0m")
		if y+2 < fb.H {
			b.WriteString("\x1b[2K\r\n")
		}
	}
	return b.String()
}

// ToANSI256 renders using the 256-color palette (fallback for terminals
// without truecolor).
func (fb *Framebuffer) ToANSI256() string {
	var b strings.Builder
	for y := 0; y < fb.H; y += 2 {
		for x := 0; x < fb.W; x++ {
			top := fb.Pix[y*fb.W+x]
			bot := top
			if y+1 < fb.H {
				bot = fb.Pix[(y+1)*fb.W+x]
			}
			fmt.Fprintf(&b, "\x1b[48;5;%dm\x1b[38;5;%dm▀\x1b[0m", xterm256(top), xterm256(bot))
		}
		if y+2 < fb.H {
			b.WriteString("\x1b[2K\r\n")
		}
	}
	return b.String()
}

func xterm256(c m32.Color) int {
	mx, mn := m32.Max(m32.Max(c.R, c.G), c.B), m32.Min(m32.Min(c.R, c.G), c.B)
	if mx-mn < 0.1 {
		v := int(c.R*24 + 0.5)
		if v < 0 {
			v = 0
		}
		if v > 23 {
			v = 23
		}
		return 232 + v
	}
	q := func(v float32) int {
		levels := [6]float32{0, 0.2, 0.4, 0.55, 0.7, 0.9}
		best, bi := float32(10), 0
		for i, l := range levels {
			if d := m32.Abs(v - l); d < best {
				best, bi = d, i
			}
		}
		return bi
	}
	return 16 + 36*q(c.R) + 6*q(c.G) + q(c.B)
}

// Print writes the ANSI frame to stdout, moving the cursor home first so the
// frame overwrites in place.
func (fb *Framebuffer) Print() {
	out := bufio.NewWriter(os.Stdout)
	out.WriteString("\x1b[H")
	out.WriteString(fb.ToANSI())
	out.Flush()
}
