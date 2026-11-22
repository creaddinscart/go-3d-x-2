// Command aa demonstrates 3d-x-2's supersampling anti-aliasing by rendering the
// same spiky scene twice: once at 1 sample (jagged) and once at 4 samples
// (smooth), writing aa_off.png and aa_on.png.
package main

import (
	"flag"
	"fmt"
	"os"

	x3d "creaddinscart.studio.go.sl/3d-x-2"
)

func buildScene() *x3d.Scene {
	sc := x3d.NewScene()
	spiky := sc.AddMesh("spiky", x3d.Icosahedron(1.1, 1),
		x3d.NewMaterial(x3d.Hex("#2b6cff")).Glossy(64, x3d.RGB(0.9, 0.9, 1)))
	spiky.Rotation = x3d.QuatFromAxis(x3d.V3(0, 1, 0), x3d.Deg2Rad(20))
	torus := sc.AddMesh("torus", x3d.Torus(1.5, 0.18, 40, 12),
		x3d.NewMaterial(x3d.Hex("#d94f30")).Textured(x3d.TexChecker(16, 4, x3d.Hex("#d94f30"), x3d.Hex("#ffd23f"))))
	torus.Rotation = x3d.QuatFromAxis(x3d.V3(1, 0, 0), x3d.Deg2Rad(70))
	return sc
}

func main() {
	off := flag.String("off", "aa_off.png", "no-AA output")
	on := flag.String("on", "aa_on.png", "AA output")
	w := flag.Int("w", 300, "render width")
	h := flag.Int("h", 200, "render height")
	flag.Parse()

	cam := x3d.NewCamera(x3d.V3(0, 0.6, 4.2), x3d.Zero3, float32(*w)/float32(*h))

	fb1 := x3d.NewFramebuffer(*w, *h)
	x3d.NewRenderer(fb1, cam).WithAA(1).Render(buildScene())
	if err := fb1.SavePNG(*off); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fb4 := x3d.NewFramebuffer(*w, *h)
	x3d.NewRenderer(fb4, cam).WithAA(4).Render(buildScene())
	if err := fb4.SavePNG(*on); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("wrote %s (1 sample) and %s (4 samples / 16x AA)\n", *off, *on)
	fmt.Println(x3d.Brand())
}
