// Command particles renders a glowing particle fountain over a glossy sphere in
// exponential fog, writing an animated GIF and a still PNG. It exercises the
// billboard/particle system, fog, additive blending and anti-aliasing.
package main

import (
	"flag"
	"fmt"
	"math"
	"os"

	x3d "creaddinscart.studio.go.sl/3d-x-2"
)

// glowSprite builds a radial alpha-falloff texture for soft particles.
func glowSprite(size int) *x3d.Texture {
	t := x3d.NewTexture(size, size)
	c := float32(size) / 2
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := (float32(x) - c) / c
			dy := (float32(y) - c) / c
			d := float32(math.Sqrt(float64(dx*dx + dy*dy)))
			a := 1 - d
			if a < 0 {
				a = 0
			}
			t.Set(x, y, x3d.NewRGBA(1, 1, 1, a*a))
		}
	}
	return t
}

func main() {
	out := flag.String("out", "fountain.gif", "output GIF path")
	png := flag.String("png", "fountain.png", "output still PNG path")
	frames := flag.Int("frames", 40, "frames")
	w := flag.Int("w", 260, "render width")
	h := flag.Int("h", 170, "render height")
	flag.Parse()

	fb := x3d.NewFramebuffer(*w, *h)
	cam := x3d.NewCamera(x3d.V3(0, 3.0, 6.5), x3d.V3(0, 0.4, -2), float32(*w)/float32(*h))

	sc := x3d.NewScene()
	sc.Background = x3d.Hex("#101826")
	sc.Fog = x3d.Fog{Mode: x3d.FogExp, Color: x3d.Hex("#101826"), Density: 0.055}
	sc.Lights = []x3d.Light{
		x3d.Ambient(x3d.Hex("#4a5668"), 1),
		x3d.Directional(x3d.Hex("#fff2e0"), 0.7, x3d.V3(-0.4, -1, -0.5)),
		x3d.PointLight(x3d.Hex("#ffaa44"), 1.1, x3d.V3(0, 2.5, 1.5), 16),
	}
	ball := sc.AddMesh("ball", x3d.Sphere(0.8, 28, 36),
		x3d.NewMaterial(x3d.Hex("#3a6bd9")).Glossy(96, x3d.RGB(1, 1, 1)))
	ball.Position = x3d.V3(0, 0.8, 0)
	sc.AddMesh("floor", x3d.Plane(60, 60, 1),
		x3d.NewMaterial(x3d.Hex("#6b7c99")).Textured(x3d.TexChecker(30, 30, x3d.Hex("#55617a"), x3d.Hex("#8497b8"))))

	ps := x3d.NewParticles(x3d.ParticleConfig{
		Origin: x3d.V3(0, 0.2, 0), Direction: x3d.V3(0, 1, 0), Speed: 4.2,
		Spread: x3d.V3(1.1, 0.6, 1.1), Rate: 140, Lifetime: 1.8,
		Size: 0.22, Color: x3d.Hex("#ffcf6b"), Gravity: 5.5,
		FadeOut: true, Additive: true, Max: 1200, Tex: glowSprite(24),
	})

	r := x3d.NewRenderer(fb, cam).WithAA(2)
	const dt = 1.0 / 30.0
	shots := make([]*x3d.Framebuffer, 0, *frames)
	for i := 0; i < *frames; i++ {
		ps.Update(dt)
		sc.Points = ps.AppendPoints(sc.Points[:0])
		r.Render(sc)
		snap := x3d.NewFramebuffer(*w, *h)
		copy(snap.Pix, fb.Pix)
		shots = append(shots, snap)
		if i == *frames-1 {
			if err := fb.SavePNG(*png); err != nil {
				fmt.Fprintln(os.Stderr, "png:", err)
			}
		}
	}

	f, err := os.Create(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()
	if err := x3d.EncodeGIF(f, shots, 3); err != nil {
		fmt.Fprintln(os.Stderr, "gif:", err)
		os.Exit(1)
	}
	fmt.Printf("wrote %s (%d frames, %d particles) and %s\n", *out, len(shots), ps.Count(), *png)
	fmt.Println(x3d.Brand())
}
