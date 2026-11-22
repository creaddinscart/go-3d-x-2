package scene

import (
	"creaddinscart.studio.go.sl/3d-x-2/geometry"
	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
	"creaddinscart.studio.go.sl/3d-x-2/texture"
)

// Entity binds geometry (or a billboard sprite) to a material.
type Entity struct {
	Mesh     *geometry.Mesh
	Material Material
}

// NewEntity builds an entity from a mesh and material.
func NewEntity(m *geometry.Mesh, mat Material) *Entity { return &Entity{Mesh: m, Material: mat} }

// Node is a transform in the scene graph.
type Node struct {
	Name     string
	Position m32.Vec3
	Rotation m32.Quaternion
	Scale    m32.Vec3
	Visible  bool

	children []*Node
	entity   *Entity
	parent   *Node
}

// NewNode creates a visible node with an identity transform.
func NewNode(name string) *Node {
	return &Node{Name: name, Scale: m32.One3, Rotation: m32.QuatIdentity(), Visible: true}
}

// AddChild attaches a child node and returns it for chaining.
func (n *Node) AddChild(c *Node) *Node { c.parent = n; n.children = append(n.children, c); return c }

// Attach binds an entity to this node.
func (n *Node) Attach(e *Entity) { n.entity = e }

// Entity returns the node's bound entity (may be nil).
func (n *Node) Entity() *Entity { return n.entity }

// Children returns the node's children.
func (n *Node) Children() []*Node { return n.children }

// Local returns the node's local transform matrix.
func (n *Node) Local() m32.Mat4 {
	return m32.Translate(n.Position).Mul(n.Rotation.Mat4()).Mul(m32.Scale(n.Scale))
}

// World composes the transform up the parent chain.
func (n *Node) World() m32.Mat4 {
	if n.parent == nil {
		return n.Local()
	}
	return n.parent.World().Mul(n.Local())
}

// Point is a camera-facing billboard used for sprites and particles. Its
// position is in world space.
type Point struct {
	Position m32.Vec3
	Color    m32.Color
	Size     float32 // world units
	Texture  *texture.Texture
	Additive bool
}

// Scene is the root of a scene graph plus its lights, background and fog.
type Scene struct {
	Root       *Node
	Lights     []Light
	Points     []Point
	Background m32.Color
	Fog        Fog
}

// New returns an empty scene with a root node and default lights.
func New() *Scene {
	return &Scene{Root: NewNode("root"), Lights: DefaultLights(), Background: m32.Hex("#0b0e14")}
}

// Add attaches a node to the scene root and returns it.
func (s *Scene) Add(n *Node) *Node { return s.Root.AddChild(n) }

// AddMesh wraps a mesh in an entity on a new node and adds it to the scene.
func (s *Scene) AddMesh(name string, mesh *geometry.Mesh, mat Material) *Node {
	n := NewNode(name)
	n.Attach(NewEntity(mesh, mat))
	s.Add(n)
	return n
}

// DrawItem is a flattened, world-transformed renderable.
type DrawItem struct {
	World m32.Mat4
	Ent   *Entity
}

// VisibleEntities walks the graph collecting renderable items.
func (s *Scene) VisibleEntities() []DrawItem {
	var out []DrawItem
	var walk func(n *Node)
	walk = func(n *Node) {
		if !n.Visible {
			return
		}
		if n.entity != nil && n.entity.Mesh != nil {
			out = append(out, DrawItem{World: n.World(), Ent: n.entity})
		}
		for _, c := range n.children {
			walk(c)
		}
	}
	walk(s.Root)
	return out
}

// Find looks up a node by name anywhere in the graph.
func (s *Scene) Find(name string) *Node {
	var found *Node
	var walk func(n *Node)
	walk = func(n *Node) {
		if found != nil {
			return
		}
		if n.Name == name {
			found = n
			return
		}
		for _, c := range n.children {
			walk(c)
		}
	}
	walk(s.Root)
	return found
}
