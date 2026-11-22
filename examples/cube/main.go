// Command cube renders a textured, lit, spinning cube with 2x supersampled
// anti-aliasing and writes a PNG (plus one ANSI terminal frame).
package main

import (
	"flag"
	"fmt"
	"os"

	x3d "creaddinscart.studio.go.sl/3d-x-2"
)

func main() {
	out := flag.String("out", "cube.png", "output PNG path")
	w := flag.Int("w", 320, "render width")
	h := flag.Int("h", 200, "render height")
	aa := flag.Int("aa", 2, "supersampling factor (1=off,2=4x,4=16x)")
	noANSI := flag.Bool("noansi", false, "skip the ANSI terminal frame")
	flag.Parse()

	fb := x3d.NewFramebuffer(*w, *h)
	cam := x3d.NewCamera(x3d.V3(0, 1.4, 4), x3d.Zero3, float32(*w)/float32(*h))

	sc := x3d.NewScene()
	cube := sc.AddMesh("cube", x3d.Cube(1.6),
		x3d.NewMaterial(x3d.Hex("#e8e8e8")).Textured(x3d.TexChecker(8, 8, x3d.Hex("#2b6cff"), x3d.White)))
	cube.Rotation = x3d.QuatFromAxis(x3d.V3(0.4, 1, 0.2), x3d.Deg2Rad(35))

	x3d.NewRenderer(fb, cam).WithAA(*aa).Render(sc)

	if err := fb.SavePNG(*out); err != nil {
		fmt.Fprintln(os.Stderr, "png:", err)
	}
	if !*noANSI {
		fmt.Print("\x1b[2J\x1b[H")
		fb.Print()
		fmt.Print("\x1b[0m\n")
	}
	fmt.Println(x3d.Brand())
}
