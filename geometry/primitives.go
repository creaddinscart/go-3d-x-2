package geometry

import (
	"math"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
)

// Cube builds a unit-centered cube of the given size with sharp face normals.
func Cube(size float32) *Mesh {
	h := size / 2
	faces := []struct {
		n          m32.Vec3
		v0, v1, v2 [3]float32
	}{
		{m32.V3(0, 0, 1), [3]float32{-h, -h, h}, [3]float32{h, -h, h}, [3]float32{h, h, h}},
		{m32.V3(0, 0, -1), [3]float32{h, -h, -h}, [3]float32{-h, -h, -h}, [3]float32{-h, h, -h}},
		{m32.V3(1, 0, 0), [3]float32{h, -h, h}, [3]float32{h, -h, -h}, [3]float32{h, h, -h}},
		{m32.V3(-1, 0, 0), [3]float32{-h, -h, -h}, [3]float32{-h, -h, h}, [3]float32{-h, h, h}},
		{m32.V3(0, 1, 0), [3]float32{-h, h, h}, [3]float32{h, h, h}, [3]float32{h, h, -h}},
		{m32.V3(0, -1, 0), [3]float32{-h, -h, -h}, [3]float32{h, -h, -h}, [3]float32{h, -h, h}},
	}
	m := New()
	uvs := [4]m32.Vec2{m32.V2(0, 0), m32.V2(1, 0), m32.V2(1, 1), m32.V2(0, 1)}
	for _, f := range faces {
		a := m32.V3(f.v0[0], f.v0[1], f.v0[2])
		b := m32.V3(f.v1[0], f.v1[1], f.v1[2])
		c := m32.V3(f.v2[0], f.v2[1], f.v2[2])
		d := a.Add(c).Sub(b)
		m.AddQuad(Vertex{a, f.n, uvs[0]}, Vertex{b, f.n, uvs[1]}, Vertex{c, f.n, uvs[2]}, Vertex{d, f.n, uvs[3]})
	}
	return m
}

// Sphere builds a UV sphere with the given radius and tessellation.
func Sphere(radius float32, stacks, slices int) *Mesh {
	if stacks < 3 {
		stacks = 3
	}
	if slices < 3 {
		slices = 3
	}
	m := New()
	for i := 0; i < stacks; i++ {
		phi0 := float32(math.Pi) * (float32(i)/float32(stacks) - 0.5)
		phi1 := float32(math.Pi) * (float32(i+1)/float32(stacks) - 0.5)
		for j := 0; j < slices; j++ {
			t0 := 2 * float32(math.Pi) * float32(j) / float32(slices)
			t1 := 2 * float32(math.Pi) * float32(j+1) / float32(slices)
			p := func(phi, theta, u, v float32) Vertex {
				n := m32.V3(float32(math.Cos(float64(phi)))*float32(math.Cos(float64(theta))),
					float32(math.Sin(float64(phi))),
					float32(math.Cos(float64(phi)))*float32(math.Sin(float64(theta))))
				return Vertex{n.Mul(radius), n, m32.V2(u, v)}
			}
			u0, u1 := float32(j)/float32(slices), float32(j+1)/float32(slices)
			v0, v1 := float32(i)/float32(stacks), float32(i+1)/float32(stacks)
			m.AddQuad(p(phi0, t0, u0, v0), p(phi0, t1, u1, v0), p(phi1, t1, u1, v1), p(phi1, t0, u0, v1))
		}
	}
	return m
}

// Plane builds an XZ plane of the given size with subdivisions.
func Plane(width, depth float32, subdiv int) *Mesh {
	if subdiv < 1 {
		subdiv = 1
	}
	m := New()
	n := m32.V3(0, 1, 0)
	for i := 0; i < subdiv; i++ {
		for j := 0; j < subdiv; j++ {
			x0 := (float32(j)/float32(subdiv) - 0.5) * width
			x1 := (float32(j+1)/float32(subdiv) - 0.5) * width
			z0 := (float32(i)/float32(subdiv) - 0.5) * depth
			z1 := (float32(i+1)/float32(subdiv) - 0.5) * depth
			u0, u1 := float32(j)/float32(subdiv), float32(j+1)/float32(subdiv)
			v0, v1 := float32(i)/float32(subdiv), float32(i+1)/float32(subdiv)
			m.AddQuad(
				Vertex{m32.V3(x0, 0, z0), n, m32.V2(u0, v0)},
				Vertex{m32.V3(x0, 0, z1), n, m32.V2(u0, v1)},
				Vertex{m32.V3(x1, 0, z1), n, m32.V2(u1, v1)},
				Vertex{m32.V3(x1, 0, z0), n, m32.V2(u1, v0)})
		}
	}
	return m
}

// Torus builds a torus with the given radii and tessellation.
func Torus(R, r float32, rings, sides int) *Mesh {
	if rings < 3 {
		rings = 3
	}
	if sides < 3 {
		sides = 3
	}
	m := New()
	at := func(i, j int) Vertex {
		u := 2 * float32(math.Pi) * float32(i) / float32(rings)
		v := 2 * float32(math.Pi) * float32(j) / float32(sides)
		cu, su := float32(math.Cos(float64(u))), float32(math.Sin(float64(u)))
		cv, sv := float32(math.Cos(float64(v))), float32(math.Sin(float64(v)))
		return Vertex{m32.V3((R+r*cv)*cu, r*sv, (R+r*cv)*su), m32.V3(cv*cu, sv, cv*su),
			m32.V2(float32(i)/float32(rings), float32(j)/float32(sides))}
	}
	for i := 0; i < rings; i++ {
		for j := 0; j < sides; j++ {
			m.AddQuad(at(i, j), at(i+1, j), at(i+1, j+1), at(i, j+1))
		}
	}
	return m
}

// Cylinder builds a capped cylinder along the Y axis.
func Cylinder(radius, height float32, segments int) *Mesh {
	if segments < 3 {
		segments = 3
	}
	m := New()
	hh := height / 2
	for i := 0; i < segments; i++ {
		a0 := 2 * float32(math.Pi) * float32(i) / float32(segments)
		a1 := 2 * float32(math.Pi) * float32(i+1) / float32(segments)
		c0, s0 := float32(math.Cos(float64(a0))), float32(math.Sin(float64(a0)))
		c1, s1 := float32(math.Cos(float64(a1))), float32(math.Sin(float64(a1)))
		n0, n1 := m32.V3(c0, 0, s0), m32.V3(c1, 0, s1)
		u0, u1 := float32(i)/float32(segments), float32(i+1)/float32(segments)
		m.AddQuad(
			Vertex{m32.V3(c0*radius, -hh, s0*radius), n0, m32.V2(u0, 0)},
			Vertex{m32.V3(c1*radius, -hh, s1*radius), n1, m32.V2(u1, 0)},
			Vertex{m32.V3(c1*radius, hh, s1*radius), n1, m32.V2(u1, 1)},
			Vertex{m32.V3(c0*radius, hh, s0*radius), n0, m32.V2(u0, 1)})
		m.AddTriangle(Vertex{m32.Zero3, m32.UnitY, m32.V2(0.5, 0.5)},
			Vertex{m32.V3(c0*radius, hh, s0*radius), m32.UnitY, m32.V2(c0*0.5+0.5, s0*0.5+0.5)},
			Vertex{m32.V3(c1*radius, hh, s1*radius), m32.UnitY, m32.V2(c1*0.5+0.5, s1*0.5+0.5)})
		m.AddTriangle(Vertex{m32.Zero3.SetY(-hh), m32.V3(0, -1, 0), m32.V2(0.5, 0.5)},
			Vertex{m32.V3(c1*radius, -hh, s1*radius), m32.V3(0, -1, 0), m32.V2(c1*0.5+0.5, s1*0.5+0.5)},
			Vertex{m32.V3(c0*radius, -hh, s0*radius), m32.V3(0, -1, 0), m32.V2(c0*0.5+0.5, s0*0.5+0.5)})
	}
	return m
}

// Cone builds a capped cone along the Y axis with its base at -height/2.
func Cone(radius, height float32, segments int) *Mesh {
	if segments < 3 {
		segments = 3
	}
	m := New()
	hh := height / 2
	apex := m32.V3(0, hh, 0)
	for i := 0; i < segments; i++ {
		a0 := 2 * float32(math.Pi) * float32(i) / float32(segments)
		a1 := 2 * float32(math.Pi) * float32(i+1) / float32(segments)
		c0, s0 := float32(math.Cos(float64(a0))), float32(math.Sin(float64(a0)))
		c1, s1 := float32(math.Cos(float64(a1))), float32(math.Sin(float64(a1)))
		b0, b1 := m32.V3(c0*radius, -hh, s0*radius), m32.V3(c1*radius, -hh, s1*radius)
		na := m32.V3(c0, radius/height, s0).Normalize()
		nb := m32.V3(c1, radius/height, s1).Normalize()
		m.AddTriangle(Vertex{b0, na, m32.V2(float32(i)/float32(segments), 0)},
			Vertex{b1, nb, m32.V2(float32(i+1)/float32(segments), 0)},
			Vertex{apex, nb, m32.V2(0.5, 1)})
		m.AddTriangle(Vertex{m32.Zero3.SetY(-hh), m32.V3(0, -1, 0), m32.V2(0.5, 0.5)},
			Vertex{b1, m32.V3(0, -1, 0), m32.V2(c1*0.5+0.5, s1*0.5+0.5)},
			Vertex{b0, m32.V3(0, -1, 0), m32.V2(c0*0.5+0.5, s0*0.5+0.5)})
	}
	return m
}

// Tetrahedron builds a regular tetrahedron inscribed in a sphere of the given
// radius, with flat face normals.
func Tetrahedron(radius float32) *Mesh {
	verts := []m32.Vec3{m32.V3(1, 1, 1), m32.V3(-1, -1, 1), m32.V3(-1, 1, -1), m32.V3(1, -1, -1)}
	for i := range verts {
		verts[i] = verts[i].Normalize().Mul(radius)
	}
	faces := [][3]int{{0, 1, 2}, {0, 3, 1}, {0, 2, 3}, {1, 3, 2}}
	m := New()
	for _, f := range faces {
		a, b, c := verts[f[0]], verts[f[1]], verts[f[2]]
		n := b.Sub(a).Cross(c.Sub(a)).Normalize()
		m.AddTriangle(Vertex{a, n, m32.V2(0, 0)}, Vertex{b, n, m32.V2(1, 0)}, Vertex{c, n, m32.V2(0.5, 1)})
	}
	return m
}

// Icosahedron builds a geodesic sphere from an icosahedron subdivided `subdiv`
// times and projected onto a sphere of the given radius.
func Icosahedron(radius float32, subdiv int) *Mesh {
	t := (1 + float32(math.Sqrt(5))) / 2
	base := []m32.Vec3{
		m32.V3(-1, t, 0), m32.V3(1, t, 0), m32.V3(-1, -t, 0), m32.V3(1, -t, 0),
		m32.V3(0, -1, t), m32.V3(0, 1, t), m32.V3(0, -1, -t), m32.V3(0, 1, -t),
		m32.V3(t, 0, -1), m32.V3(t, 0, 1), m32.V3(-t, 0, -1), m32.V3(-t, 0, 1),
	}
	for i := range base {
		base[i] = base[i].Normalize()
	}
	faces := [][3]int{
		{0, 11, 5}, {0, 5, 1}, {0, 1, 7}, {0, 7, 10}, {0, 10, 11},
		{1, 5, 9}, {5, 11, 4}, {11, 10, 2}, {10, 7, 6}, {7, 1, 8},
		{3, 9, 4}, {3, 4, 2}, {3, 2, 6}, {3, 6, 8}, {3, 8, 9},
		{4, 9, 5}, {2, 4, 11}, {6, 2, 10}, {8, 6, 7}, {9, 8, 1},
	}
	cur := make([][3]m32.Vec3, 0, len(faces))
	for _, f := range faces {
		cur = append(cur, [3]m32.Vec3{base[f[0]], base[f[1]], base[f[2]]})
	}
	for s := 0; s < subdiv; s++ {
		var next [][3]m32.Vec3
		for _, tri := range cur {
			ab := tri[0].Add(tri[1]).Normalize()
			bc := tri[1].Add(tri[2]).Normalize()
			ca := tri[2].Add(tri[0]).Normalize()
			next = append(next,
				[3]m32.Vec3{tri[0], ab, ca}, [3]m32.Vec3{tri[1], bc, ab},
				[3]m32.Vec3{tri[2], ca, bc}, [3]m32.Vec3{ab, bc, ca})
		}
		cur = next
	}
	m := New()
	for _, tri := range cur {
		a, b, c := tri[0].Normalize(), tri[1].Normalize(), tri[2].Normalize()
		n := b.Sub(a).Cross(c.Sub(a)).Normalize()
		m.AddTriangle(Vertex{a.Mul(radius), n, m32.V2(0, 0)},
			Vertex{b.Mul(radius), n, m32.V2(1, 0)}, Vertex{c.Mul(radius), n, m32.V2(0.5, 1)})
	}
	return m
}
