package render

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
)

// webSafePalette is the 6x6x6 color cube (216 entries) used to quantize GIF
// frames without a full median-cut.
var webSafePalette = buildWebSafePalette()

func buildWebSafePalette() []color.Color {
	levels := []uint8{0, 51, 102, 153, 204, 255}
	pal := make([]color.Color, 0, 216)
	for _, r := range levels {
		for _, g := range levels {
			for _, b := range levels {
				pal = append(pal, color.NRGBA{R: r, G: g, B: b, A: 255})
			}
		}
	}
	return pal
}

func paletteIndex(c m32.Color) uint8 {
	q := func(v float32) int {
		n := int(v*5 + 0.5)
		if n < 0 {
			return 0
		}
		if n > 5 {
			return 5
		}
		return n
	}
	return uint8(36*q(c.R) + 6*q(c.G) + q(c.B))
}

// EncodeGIF writes the frames as an animated GIF. delayCS is the per-frame
// delay in centiseconds (1/100 s).
func EncodeGIF(w io.Writer, frames []*Framebuffer, delayCS int) error {
	if len(frames) == 0 {
		return nil
	}
	g := &gif.GIF{LoopCount: 0}
	for _, fb := range frames {
		img := image.NewPaletted(image.Rect(0, 0, fb.W, fb.H), webSafePalette)
		for i := 0; i < fb.W*fb.H; i++ {
			img.Pix[i] = paletteIndex(fb.Pix[i])
		}
		g.Image = append(g.Image, img)
		g.Delay = append(g.Delay, delayCS)
	}
	return gif.EncodeAll(w, g)
}

// SaveGIF writes the frames to a GIF file on disk.
func SaveGIF(path string, frames []*Framebuffer, delayCS int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return EncodeGIF(f, frames, delayCS)
}
