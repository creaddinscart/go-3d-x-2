// Package scene defines the scene graph, materials, lights, camera and fog for
// 3d-x-2. It is consumed by the render package.
//
// Built by CREADDINSCART STUDIO — https://www.shit.pub/studio/
package scene

import (
	"math"

	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
	"creaddinscart.studio.go.sl/3d-x-2/texture"
)

// ShadeMode selects the per-pixel interpolation strategy.
type ShadeMode int

const (
	ShadeFlat ShadeMode = iota
	ShadeGouraud
	ShadePhong
)

// Material describes how a surface reflects light.
type Material struct {
	Ambient     m32.Color
	Diffuse     m32.Color
	Specular    m32.Color
	Shininess   float32
	Emissive    m32.Color
	Texture     *texture.Texture
	ShadeMode   ShadeMode
	DoubleSided bool
	Opacity     float32 // <1 enables alpha blending
}

// NewMaterial returns a Phong-shaded matte material with the given base color.
func NewMaterial(base m32.Color) Material {
	return Material{
		Ambient:   base.Mul(0.35),
		Diffuse:   base,
		Specular:  m32.RGB(0.2, 0.2, 0.2),
		Shininess: 16,
		Emissive:  m32.Black,
		ShadeMode: ShadePhong,
		Opacity:   1,
	}
}

// Textured returns a copy using the given texture as the diffuse map.
func (m Material) Textured(t *texture.Texture) Material { m.Texture = t; return m }

// Glossy returns a copy with a stronger specular response.
func (m Material) Glossy(shininess float32, spec m32.Color) Material {
	m.Shininess = shininess
	m.Specular = spec
	return m
}

// Flat returns a copy using flat shading.
func (m Material) Flat() Material { m.ShadeMode = ShadeFlat; return m }

// Transparent returns a copy with the given opacity (alpha-blended).
func (m Material) Transparent(opacity float32) Material {
	m.Opacity = m32.Clamp01(opacity)
	return m
}

// LightKind selects a lighting model.
type LightKind int

const (
	LightAmbient LightKind = iota
	LightDirectional
	LightPoint
	LightSpot
)

// Light is a single scene light in world space.
type Light struct {
	Kind      LightKind
	Color     m32.Color
	Intensity float32
	Direction m32.Vec3
	Position  m32.Vec3
	Range     float32
	SpotCos   float32
}

// Ambient builds an ambient light.
func Ambient(c m32.Color, intensity float32) Light {
	return Light{Kind: LightAmbient, Color: c, Intensity: intensity}
}

// Directional builds a directional light pointing along dir.
func Directional(c m32.Color, intensity float32, dir m32.Vec3) Light {
	return Light{Kind: LightDirectional, Color: c, Intensity: intensity, Direction: dir.Normalize()}
}

// PointLight builds a point light at pos with optional range (0 = no falloff).
func PointLight(c m32.Color, intensity float32, pos m32.Vec3, rng float32) Light {
	return Light{Kind: LightPoint, Color: c, Intensity: intensity, Position: pos, Range: rng}
}

// Spot builds a spot light at pos aimed along dir with a cone half-angle.
func Spot(c m32.Color, intensity float32, pos, dir m32.Vec3, angleRad, rng float32) Light {
	return Light{
		Kind: LightSpot, Color: c, Intensity: intensity,
		Position: pos, Direction: dir.Normalize(), Range: rng,
		SpotCos: float32(math.Cos(float64(angleRad))),
	}
}

// DefaultLights returns a pleasant three-light studio rig.
func DefaultLights() []Light {
	return []Light{
		Ambient(m32.Hex("#2e2e38"), 1),
		Directional(m32.Hex("#fff5e8"), 0.9, m32.V3(-0.5, -1, -0.35)),
		PointLight(m32.Hex("#66a3ff"), 0.6, m32.V3(3, 3, 3), 12),
	}
}

// FogMode selects the depth-fog falloff curve.
type FogMode int

const (
	FogNone FogMode = iota
	FogLinear
	FogExp
	FogExp2
)

// Fog blends distant fragments toward a color to fake atmosphere.
type Fog struct {
	Mode    FogMode
	Color   m32.Color
	Start   float32 // linear
	End     float32 // linear
	Density float32 // exp / exp2
}

// Factor returns the fog blend weight [0,1] for a view-space depth.
func (f Fog) Factor(depth float32) float32 {
	if f.Mode == FogNone {
		return 0
	}
	switch f.Mode {
	case FogLinear:
		return m32.Clamp01((depth - f.Start) / (f.End - f.Start))
	case FogExp:
		return 1 - float32(math.Exp(-float64(f.Density*depth)))
	case FogExp2:
		x := f.Density * depth
		return 1 - float32(math.Exp(-float64(x*x)))
	}
	return 0
}
