package render

import (
	"math"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
	"creaddinscart.studio.go.sl/3d-x-2/scene"
)

type clipVertex struct {
	clip   m32.Vec4
	world  m32.Vec3
	normal m32.Vec3
	uv     m32.Vec2
}

type screenVertex struct {
	x, y, z float32
	invW    float32
	world   m32.Vec3
	normal  m32.Vec3
	uv      m32.Vec2
	col     m32.Color
}

// frustum planes in clip space: keep points where a*x+b*y+c*z+d*w >= 0.
var frustumPlanes = [6]m32.Vec4{
	m32.V4(1, 0, 0, 1), m32.V4(-1, 0, 0, 1),
	m32.V4(0, 1, 0, 1), m32.V4(0, -1, 0, 1),
	m32.V4(0, 0, 1, 1), m32.V4(0, 0, -1, 1),
}

func planeDist(p, pl m32.Vec4) float32 {
	return pl.X*p.X + pl.Y*p.Y + pl.Z*p.Z + pl.W*p.W
}

func lerpCV(a, b clipVertex, t float32) clipVertex {
	return clipVertex{
		clip:   a.clip.Mul(1 - t).Add(b.clip.Mul(t)),
		world:  a.world.Lerp(b.world, t),
		normal: a.normal.Lerp(b.normal, t).Normalize(),
		uv:     m32.V2(m32.Lerp(a.uv.X, b.uv.X, t), m32.Lerp(a.uv.Y, b.uv.Y, t)),
	}
}

// clipTriangle clips against all frustum planes, returning a convex polygon.
func clipTriangle(tri [3]clipVertex) []clipVertex {
	poly := tri[:]
	for _, pl := range frustumPlanes {
		if len(poly) == 0 {
			return nil
		}
		var out []clipVertex
		for i := 0; i < len(poly); i++ {
			cur, nxt := poly[i], poly[(i+1)%len(poly)]
			dc := planeDist(cur.clip, pl)
			dn := planeDist(nxt.clip, pl)
			if dc >= 0 {
				out = append(out, cur)
			}
			if (dc >= 0) != (dn >= 0) {
				out = append(out, lerpCV(cur, nxt, dc/(dc-dn)))
			}
		}
		poly = out
	}
	return poly
}

func (r *Renderer) drawMesh(e *scene.Entity, mvp, world m32.Mat4) {
	m := e.Mesh
	vc := make([]clipVertex, len(m.Vertices))
	for i, v := range m.Vertices {
		vc[i] = clipVertex{
			clip:   mvp.MulVec4(v.Position.ToVec4(1)),
			world:  world.TransformPoint(v.Position),
			normal: world.TransformVector(v.Normal).Normalize(),
			uv:     v.UV,
		}
	}
	for t := 0; t+2 < len(m.Indices); t += 3 {
		tri := [3]clipVertex{vc[m.Indices[t]], vc[m.Indices[t+1]], vc[m.Indices[t+2]]}
		poly := clipTriangle(tri)
		if len(poly) < 3 {
			continue
		}
		if r.Wireframe {
			for i := 0; i < len(poly); i++ {
				a := r.project(poly[i], e.Material)
				b := r.project(poly[(i+1)%len(poly)], e.Material)
				r.drawEdge(a, b)
			}
			continue
		}
		for i := 1; i+1 < len(poly); i++ {
			r.rasterize(poly[0], poly[i], poly[i+1], e.Material)
		}
	}
}

func (r *Renderer) project(cv clipVertex, mat scene.Material) screenVertex {
	w := cv.clip.W
	if w == 0 {
		w = 1e-6
	}
	ndc := m32.V3(cv.clip.X/w, cv.clip.Y/w, cv.clip.Z/w)
	sv := screenVertex{
		x: (ndc.X*0.5 + 0.5) * float32(r.cur.W),
		y: (1 - (ndc.Y*0.5 + 0.5)) * float32(r.cur.H),
		z: ndc.Z, invW: 1 / w, world: cv.world, normal: cv.normal, uv: cv.uv,
	}
	if mat.ShadeMode == scene.ShadeFlat || mat.ShadeMode == scene.ShadeGouraud {
		sv.col = r.shade(mat, cv.world, cv.normal, cv.uv)
	}
	return sv
}

// rasterize uses incremental edge functions: the three half-plane values are
// computed once at the top-left pixel and then stepped by integer deltas, so
// the inner loop is a handful of adds.
func (r *Renderer) rasterize(a, b, c clipVertex, mat scene.Material) {
	sa, sb, sc := r.project(a, mat), r.project(b, mat), r.project(c, mat)

	// edge coefficients A,B (C folded into the start value)
	aA, bA := sb.y-sc.y, sc.x-sb.x
	aB, bB := sc.y-sa.y, sa.x-sc.x
	aC, bC := sa.y-sb.y, sb.x-sa.x
	area := aA*sa.x + bA*sa.y + (sb.x*sc.y - sc.x*sb.y) // 2*signed area
	if area == 0 {
		return
	}
	if r.CullBackface && !mat.DoubleSided && area > 0 {
		return
	}

	if mat.ShadeMode == scene.ShadeFlat {
		faceN := sb.world.Sub(sa.world).Cross(sc.world.Sub(sa.world)).Normalize()
		if faceN.Dot(sa.normal) < 0 {
			faceN = faceN.Neg()
		}
		col := r.shade(mat, sa.world, faceN, sa.uv)
		sa.col, sb.col, sc.col = col, col, col
	}

	minX := int(math.Floor(float64(m32.Min3(sa.x, sb.x, sc.x))))
	maxX := int(math.Ceil(float64(m32.Max3(sa.x, sb.x, sc.x))))
	minY := int(math.Floor(float64(m32.Min3(sa.y, sb.y, sc.y))))
	maxY := int(math.Ceil(float64(m32.Max3(sa.y, sb.y, sc.y))))
	minX = int(m32.Max(float32(minX), 0))
	minY = int(m32.Max(float32(minY), 0))
	maxX = int(m32.Min(float32(maxX), float32(r.cur.W-1)))
	maxY = int(m32.Min(float32(maxY), float32(r.cur.H-1)))

	invArea := 1 / area
	const cEdge = 0.5
	px0 := float32(minX) + cEdge
	cA := sb.x*sc.y - sc.x*sb.y
	cB := sc.x*sa.y - sa.x*sc.y
	cC := sa.x*sb.y - sb.x*sa.y

	for py := minY; py <= maxY; py++ {
		y := float32(py) + cEdge
		w0 := aA*px0 + bA*y + cA
		w1 := aB*px0 + bB*y + cB
		w2 := aC*px0 + bC*y + cC
		for px := minX; px <= maxX; px++ {
			l0, l1, l2 := w0*invArea, w1*invArea, w2*invArea
			if l0 >= 0 && l1 >= 0 && l2 >= 0 {
				iw := l0*sa.invW + l1*sb.invW + l2*sc.invW
				k := 1 / iw
				z := l0*sa.z + l1*sb.z + l2*sc.z
				idx := py*r.cur.W + px
				if z < r.cur.Depth[idx] {
					r.cur.Depth[idx] = z
					r.cur.Pix[idx] = r.shadeInterp(mat, l0, l1, l2, sa, sb, sc, k)
				}
			}
			w0 += aA
			w1 += aB
			w2 += aC
		}
	}
}

func (r *Renderer) shadeInterp(mat scene.Material, l0, l1, l2 float32, sa, sb, sc screenVertex, k float32) m32.Color {
	if mat.ShadeMode == scene.ShadePhong {
		wn := m32.V3(
			(l0*sa.invW)*sa.normal.X+(l1*sb.invW)*sb.normal.X+(l2*sc.invW)*sc.normal.X,
			(l0*sa.invW)*sa.normal.Y+(l1*sb.invW)*sb.normal.Y+(l2*sc.invW)*sc.normal.Y,
			(l0*sa.invW)*sa.normal.Z+(l1*sb.invW)*sb.normal.Z+(l2*sc.invW)*sc.normal.Z,
		).Mul(k).Normalize()
		wp := m32.V3(
			(l0*sa.invW)*sa.world.X+(l1*sb.invW)*sb.world.X+(l2*sc.invW)*sc.world.X,
			(l0*sa.invW)*sa.world.Y+(l1*sb.invW)*sb.world.Y+(l2*sc.invW)*sc.world.Y,
			(l0*sa.invW)*sa.world.Z+(l1*sb.invW)*sb.world.Z+(l2*sc.invW)*sc.world.Z,
		).Mul(k)
		uv := m32.V2(
			(l0*sa.invW)*sa.uv.X+(l1*sb.invW)*sb.uv.X+(l2*sc.invW)*sc.uv.X,
			(l0*sa.invW)*sa.uv.Y+(l1*sb.invW)*sb.uv.Y+(l2*sc.invW)*sc.uv.Y,
		).Mul(k)
		return r.shade(mat, wp, wn, uv)
	}
	cR, cG, cB := l0*sa.invW, l1*sb.invW, l2*sc.invW
	s := cR + cG + cB
	if s == 0 {
		return m32.NewRGBA(0, 0, 0, mat.Opacity)
	}
	return m32.NewRGBA(
		(cR*sa.col.R+cG*sb.col.R+cB*sc.col.R)/s,
		(cR*sa.col.G+cG*sb.col.G+cB*sc.col.G)/s,
		(cR*sa.col.B+cG*sb.col.B+cB*sc.col.B)/s,
		mat.Opacity,
	)
}

// drawEdge rasterizes a wireframe line with simple depth from the average z.
func (r *Renderer) drawEdge(a, b screenVertex) {
	col := m32.NewRGBA(1, 1, 1, 1)
	if a.col != (m32.Color{}) {
		col = a.col
	}
	x0, y0 := int(a.x), int(a.y)
	x1, y1 := int(b.x), int(b.y)
	dx, dy := absi(x1-x0), -absi(y1-y0)
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	err := dx + dy
	z := (a.z + b.z) * 0.5
	for {
		if x0 >= 0 && x0 < r.cur.W && y0 >= 0 && y0 < r.cur.H {
			idx := y0*r.cur.W + x0
			if z < r.cur.Depth[idx] {
				r.cur.Depth[idx] = z
				r.cur.Pix[idx] = col
			}
		}
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 >= dy {
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

func absi(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// drawPoints renders scene points as camera-facing billboards (sprites and
// particles) with depth test plus alpha or additive blending.
func (r *Renderer) drawPoints(vp m32.Mat4, s *scene.Scene) {
	if len(s.Points) == 0 {
		return
	}
	// camera right in world space (first row of the view matrix)
	view := r.Camera.View()
	right := m32.V3(view.At(0, 0), view.At(1, 0), view.At(2, 0))

	for _, p := range s.Points {
		c := vp.MulVec4(p.Position.ToVec4(1))
		if c.W <= 0 {
			continue
		}
		ndc := m32.V3(c.X/c.W, c.Y/c.W, c.Z/c.W)
		cx := (ndc.X*0.5 + 0.5) * float32(r.cur.W)
		cy := (1 - (ndc.Y*0.5 + 0.5)) * float32(r.cur.H)
		off := vp.MulVec4(p.Position.Add(right.Mul(p.Size * 0.5)).ToVec4(1))
		if off.W <= 0 {
			continue
		}
		ox := (off.X/off.W*0.5 + 0.5) * float32(r.cur.W)
		rad := m32.Abs(ox - cx)
		if rad < 0.5 {
			rad = 0.5
		}
		r.drawBillboard(int(cx), int(cy), int(rad+0.5), ndc.Z, p)
	}
}

func (r *Renderer) drawBillboard(cx, cy, rad int, depth float32, p scene.Point) {
	for dy := -rad; dy <= rad; dy++ {
		for dx := -rad; dx <= rad; dx++ {
			if dx*dx+dy*dy > rad*rad {
				continue
			}
			px, py := cx+dx, cy+dy
			if px < 0 || px >= r.cur.W || py < 0 || py >= r.cur.H {
				continue
			}
			idx := py*r.cur.W + px
			if depth >= r.cur.Depth[idx] {
				continue
			}
			col := p.Color
			if p.Texture != nil {
				u := (float32(dx)/float32(rad) + 1) * 0.5
				v := (float32(dy)/float32(rad) + 1) * 0.5
				col = col.Modulate(p.Texture.Sample(u, 1-v))
			}
			if p.Additive {
				r.cur.Add(px, py, col.Mul(col.A))
			} else {
				r.cur.Blend(px, py, col)
			}
			r.cur.Depth[idx] = depth
		}
	}
}
