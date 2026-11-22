// Package geometry provides the mesh type, parametric primitives and the OBJ
// importer/exporter for 3d-x-2.
//
// Built by CREADDINSCART STUDIO — https://www.shit.pub/studio/
package geometry

import "creaddinscart.studio.go.sl/3d-x-2/math32"

// Vertex is a single rasterizable vertex.
type Vertex struct {
	Position math32.Vec3
	Normal   math32.Vec3
	UV       math32.Vec2
}

// Mesh is an indexed triangle list.
type Mesh struct {
	Vertices []Vertex
	Indices  []uint32
}

// New returns an empty mesh.
func New() *Mesh { return &Mesh{} }

// AddTriangle appends three vertices and their indices.
func (m *Mesh) AddTriangle(a, b, c Vertex) {
	i := uint32(len(m.Vertices))
	m.Vertices = append(m.Vertices, a, b, c)
	m.Indices = append(m.Indices, i, i+1, i+2)
}

// AddQuad appends a quad as two triangles (a,b,c,d in order).
func (m *Mesh) AddQuad(a, b, c, d Vertex) {
	m.AddTriangle(a, b, c)
	m.AddTriangle(a, c, d)
}

// TriangleCount returns the number of triangles in the mesh.
func (m *Mesh) TriangleCount() int { return len(m.Indices) / 3 }

// Bounds returns the axis-aligned bounding box of the vertex positions.
func (m *Mesh) Bounds() (min, max math32.Vec3) {
	if len(m.Vertices) == 0 {
		return
	}
	min, max = m.Vertices[0].Position, m.Vertices[0].Position
	for _, v := range m.Vertices {
		p := v.Position
		min.X, max.X = math32.Min(p.X, min.X), math32.Max(p.X, max.X)
		min.Y, max.Y = math32.Min(p.Y, min.Y), math32.Max(p.Y, max.Y)
		min.Z, max.Z = math32.Min(p.Z, min.Z), math32.Max(p.Z, max.Z)
	}
	return
}

// ComputeNormals recomputes smooth (averaged) vertex normals from faces.
func (m *Mesh) ComputeNormals() {
	for i := range m.Vertices {
		m.Vertices[i].Normal = math32.Zero3
	}
	for i := 0; i+2 < len(m.Indices); i += 3 {
		a := &m.Vertices[m.Indices[i]]
		b := &m.Vertices[m.Indices[i+1]]
		c := &m.Vertices[m.Indices[i+2]]
		n := b.Position.Sub(a.Position).Cross(c.Position.Sub(a.Position))
		a.Normal = a.Normal.Add(n)
		b.Normal = b.Normal.Add(n)
		c.Normal = c.Normal.Add(n)
	}
	for i := range m.Vertices {
		m.Vertices[i].Normal = m.Vertices[i].Normal.Normalize()
	}
}

// Transform returns a new mesh with a matrix applied to positions and normals.
func (m *Mesh) Transform(mat math32.Mat4) *Mesh {
	out := &Mesh{Vertices: make([]Vertex, len(m.Vertices)), Indices: append([]uint32(nil), m.Indices...)}
	for i, v := range m.Vertices {
		nv := v
		nv.Position = mat.TransformPoint(v.Position)
		nv.Normal = mat.TransformVector(v.Normal).Normalize()
		out.Vertices[i] = nv
	}
	return out
}
