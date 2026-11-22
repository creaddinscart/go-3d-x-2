package x3d

import (
	"bytes"
	"testing"
)

func approx(a, b, eps float32) bool { return Abs(a-b) < eps }

func TestMat4Inverse(t *testing.T) {
	m := Translate(V3(1, 2, 3)).Mul(RotateY(Deg2Rad(37))).Mul(Scale(V3(2, 3, 4)))
	id := m.Mul(m.Inverse())
	for i := 0; i < 16; i++ {
		want := float32(0)
		if i%5 == 0 {
			want = 1
		}
		if !approx(id[i], want, 1e-3) {
			t.Fatalf("inverse product at %d = %v, want %v", i, id[i], want)
		}
	}
}

func TestQuaternionRoundTrip(t *testing.T) {
	q := QuatFromAxis(V3(1, 2, 3), Deg2Rad(50))
	p := q.Mat4().TransformPoint(UnitX)
	if !approx(p.Length(), 1, 1e-3) {
		t.Fatalf("rotated length = %v, want 1", p.Length())
	}
}

func TestColorHex(t *testing.T) {
	c := Hex("#ff8000")
	if !approx(c.R, 1, 1e-3) || !approx(c.G, 128.0/255, 1e-2) || !approx(c.B, 0, 1e-3) {
		t.Fatalf("Hex parse wrong: %+v", c)
	}
}

func TestTextureSample(t *testing.T) {
	ck := TexChecker(2, 2, Black, White)
	if ck.Sample(0.1, 0.1).R == ck.Sample(0.6, 0.1).R {
		t.Fatal("checker did not alternate")
	}
}

func TestEasing(t *testing.T) {
	if !approx(EaseOutQuad(0.5), 0.75, 1e-5) {
		t.Fatalf("EaseOutQuad(0.5) = %v", EaseOutQuad(0.5))
	}
}

func TestFogFactor(t *testing.T) {
	lin := Fog{Mode: FogLinear, Start: 2, End: 10}
	if !approx(lin.Factor(2), 0, 1e-5) || !approx(lin.Factor(10), 1, 1e-5) || !approx(lin.Factor(6), 0.5, 1e-5) {
		t.Fatalf("linear fog wrong: %v %v %v", lin.Factor(2), lin.Factor(6), lin.Factor(10))
	}
	exp := Fog{Mode: FogExp, Density: 0.1}
	if !(exp.Factor(0) >= 0 && exp.Factor(50) > exp.Factor(5)) {
		t.Fatal("exp fog not monotonic")
	}
}

func TestOBJRoundTrip(t *testing.T) {
	cube := Cube(1)
	var buf bytes.Buffer
	if err := cube.WriteOBJ(&buf); err != nil {
		t.Fatal(err)
	}
	back, err := LoadOBJ(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if back.TriangleCount() != cube.TriangleCount() {
		t.Fatalf("triangles %d != %d", back.TriangleCount(), cube.TriangleCount())
	}
}

func TestRasterizerDrawsSilhouette(t *testing.T) {
	fb := NewFramebuffer(120, 80)
	cam := NewCamera(V3(0, 0, 4), Zero3, 120.0/80.0)
	sc := NewScene()
	sc.AddMesh("cube", Cube(2), NewMaterial(Red))
	NewRenderer(fb, cam).Render(sc)
	nonBg := 0
	for _, p := range fb.Pix {
		if p != sc.Background {
			nonBg++
		}
	}
	if float32(nonBg) < float32(len(fb.Pix))*0.1 {
		t.Fatalf("rasterizer drew too few pixels: %d/%d", nonBg, len(fb.Pix))
	}
}

func TestAASmoothsEdges(t *testing.T) {
	// With AA on, an edge pixel should be a blend (not exactly background nor
	// exactly object color). We just assert the AA render fills and differs.
	cam := NewCamera(V3(0, 0, 4), Zero3, 1)
	sc := NewScene()
	sc.AddMesh("s", Sphere(1.5, 20, 24), NewMaterial(White))

	off := NewFramebuffer(80, 80)
	NewRenderer(off, cam).WithAA(1).Render(sc)
	on := NewFramebuffer(80, 80)
	NewRenderer(on, cam).WithAA(2).Render(sc)

	diff := 0
	for i := range off.Pix {
		if off.Pix[i] != on.Pix[i] {
			diff++
		}
	}
	if diff == 0 {
		t.Fatal("AA produced no difference vs no-AA")
	}
}

func TestAnimationAppliesPosition(t *testing.T) {
	sc := NewScene()
	n := sc.AddMesh("m", Cube(1), NewMaterial(White))
	clip := NewClip("slide")
	clip.AddPosition("m", []Keyframe{
		{Time: 0, Value: V3(-2, 0, 0)},
		{Time: 1, Value: V3(2, 0, 0), Ease: EaseLinear},
	})
	an := NewAnimator(clip, sc)
	an.Advance(0.5)
	if !approx(n.Position.X, 0, 1e-3) {
		t.Fatalf("animated X = %v, want 0", n.Position.X)
	}
}

func TestParticles(t *testing.T) {
	ps := NewParticles(ParticleConfig{
		Origin: Zero3, Direction: UnitY, Speed: 3, Rate: 100, Lifetime: 1,
		Size: 0.2, Color: White, Max: 500,
	})
	for i := 0; i < 10; i++ {
		ps.Update(0.05)
	}
	if ps.Count() == 0 {
		t.Fatal("no particles spawned")
	}
	pts := ps.AppendPoints(nil)
	if len(pts) != ps.Count() {
		t.Fatalf("points %d != particles %d", len(pts), ps.Count())
	}
}

func TestBillboardRenders(t *testing.T) {
	fb := NewFramebuffer(80, 80)
	cam := NewCamera(V3(0, 0, 4), Zero3, 1)
	sc := NewScene()
	sc.Points = []Point{{Position: Zero3, Color: Red, Size: 2}}
	NewRenderer(fb, cam).Render(sc)
	nonBg := 0
	for _, p := range fb.Pix {
		if p != sc.Background {
			nonBg++
		}
	}
	if nonBg == 0 {
		t.Fatal("billboard point drew nothing")
	}
}

func TestFramebufferANSI(t *testing.T) {
	fb := NewFramebuffer(4, 4)
	fb.Fill(RGB(1, 0, 0))
	s := fb.ToANSI()
	if !bytes.Contains([]byte(s), []byte("▀")) || !bytes.Contains([]byte(s), []byte("48;2;")) {
		t.Fatal("ANSI output malformed")
	}
}

func TestPlaneFrontFacingFromAbove(t *testing.T) {
	fb := NewFramebuffer(60, 60)
	cam := NewCamera(V3(0, 5, 0.001), Zero3, 1) // looking straight down
	cam.Up = V3(0, 0, -1)
	sc := NewScene()
	sc.AddMesh("floor", Plane(10, 10, 1), NewMaterial(White)) // +Y normal
	NewRenderer(fb, cam).Render(sc)
	nonBg := 0
	for _, p := range fb.Pix {
		if p != sc.Background {
			nonBg++
		}
	}
	if float32(nonBg) < float32(len(fb.Pix))*0.5 {
		t.Fatalf("plane not visible from above (culled?): %d/%d", nonBg, len(fb.Pix))
	}
}

func BenchmarkRenderSphere(b *testing.B) {
	fb := NewFramebuffer(200, 120)
	cam := NewCamera(V3(0, 0, 4), Zero3, 200.0/120.0)
	sc := NewScene()
	sc.AddMesh("sphere", Sphere(1.5, 24, 32), NewMaterial(White))
	r := NewRenderer(fb, cam)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Render(sc)
	}
}
