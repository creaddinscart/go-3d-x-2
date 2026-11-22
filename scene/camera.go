package scene

import "math"

import m32 "creaddinscart.studio.go.sl/3d-x-2/math32"

// Camera holds view and projection state.
type Camera struct {
	Position m32.Vec3
	Target   m32.Vec3
	Up       m32.Vec3

	FOVY         float32 // radians (perspective)
	Aspect       float32
	Near         float32
	Far          float32
	Orthographic bool
	OrthoHeight  float32
}

// New builds a perspective camera looking from eye toward target.
func NewCamera(eye, target m32.Vec3, aspect float32) *Camera {
	return &Camera{
		Position: eye, Target: target, Up: m32.UnitY,
		FOVY: float32(60 * math.Pi / 180), Aspect: aspect,
		Near: 0.05, Far: 1000,
	}
}

// View returns the world-to-view matrix.
func (c *Camera) View() m32.Mat4 { return m32.LookAt(c.Position, c.Target, c.Up) }

// Projection returns the view-to-clip matrix.
func (c *Camera) Projection() m32.Mat4 {
	if c.Orthographic {
		h := c.OrthoHeight / 2
		w := h * c.Aspect
		return m32.Ortho(-w, w, -h, h, c.Near, c.Far)
	}
	return m32.Perspective(c.FOVY, c.Aspect, c.Near, c.Far)
}

// ViewProjection returns projection * view.
func (c *Camera) ViewProjection() m32.Mat4 { return c.Projection().Mul(c.View()) }

// Orbit returns the eye position for a spherical orbit around Target.
func (c *Camera) Orbit(yaw, pitch, dist float32) m32.Vec3 {
	x := dist * float32(math.Cos(float64(pitch))) * float32(math.Sin(float64(yaw)))
	y := dist * float32(math.Sin(float64(pitch)))
	z := dist * float32(math.Cos(float64(pitch))) * float32(math.Cos(float64(yaw)))
	return c.Target.Add(m32.V3(x, y, z))
}

// LookDir returns the normalized forward direction of the camera.
func (c *Camera) LookDir() m32.Vec3 { return c.Target.Sub(c.Position).Normalize() }

// ScreenRay casts a ray through normalized device coordinates in [-1, 1].
func (c *Camera) ScreenRay(ndc m32.Vec2) (origin, dir m32.Vec3) {
	origin = c.Position
	vp := c.ViewProjection().Inverse()
	p := vp.MulVec4(m32.V4(ndc.X, ndc.Y, 1, 1)).XYZ()
	q := vp.MulVec4(m32.V4(ndc.X, ndc.Y, -1, 1)).XYZ()
	return origin, p.Sub(q).Normalize()
}
