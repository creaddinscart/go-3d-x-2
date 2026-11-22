// Package x3d is the convenience façade of 3d-x-2, an enhanced dependency-free
// 3D graphics library for Go. It re-exports the core types and constructors of
// the sub-packages (math32, geometry, texture, scene, render, anim) so that a
// single import is enough for everyday use:
//
//	import x3d "creaddinscart.studio.go.sl/3d-x-2"
//
//	fb  := x3d.NewFramebuffer(320, 200)
//	cam := x3d.NewCamera(x3d.V3(0, 1, 4), x3d.Zero3, 1.6)
//	sc  := x3d.NewScene()
//	sc.AddMesh("cube", x3d.Cube(1.5), x3d.NewMaterial(x3d.Hex("#2b6cff")))
//	x3d.NewRenderer(fb, cam).WithAA(2).Render(sc) // 2x supersampled anti-aliasing
//	fb.SavePNG("cube.png")
//
// 3d-x-2 adds, over 3d-x: a faster incremental edge-function rasterizer,
// supersampling anti-aliasing, depth fog, wireframe mode, camera-facing
// billboards/sprites, a GPU-free particle system, and a clearer sub-package
// structure. Standard library only — no GPU, no cgo, no third-party modules.
//
// Built by CREADDINSCART STUDIO — Made with CREADDINSCART STUDIO.
// https://www.shit.pub/studio/
package x3d

import (
	"creaddinscart.studio.go.sl/3d-x-2/anim"
	"creaddinscart.studio.go.sl/3d-x-2/geometry"
	m32 "creaddinscart.studio.go.sl/3d-x-2/math32"
	"creaddinscart.studio.go.sl/3d-x-2/render"
	"creaddinscart.studio.go.sl/3d-x-2/scene"
	"creaddinscart.studio.go.sl/3d-x-2/texture"
)

// Version is the released version of the library.
const Version = "2.0.0"

// Brand returns the attribution required by the CREADDINSCART STUDIO licence.
func Brand() string {
	return "Built by CREADDINSCART STUDIO — Made with CREADDINSCART STUDIO — https://www.shit.pub/studio/"
}

// --- math32 ---------------------------------------------------------------

type (
	Vec2       = m32.Vec2
	Vec3       = m32.Vec3
	Vec4       = m32.Vec4
	Mat4       = m32.Mat4
	Quaternion = m32.Quaternion
	Color      = m32.Color
)

var (
	V2, V3, V4 = m32.V2, m32.V3, m32.V4
	RGB        = m32.RGB
	NewRGBA    = m32.NewRGBA
	Hex        = m32.Hex
	HSV        = m32.HSV

	Zero3, UnitX, UnitY, UnitZ, One3 = m32.Zero3, m32.UnitX, m32.UnitY, m32.UnitZ, m32.One3
	Black, White, Red, Green, Blue   = m32.Black, m32.White, m32.Red, m32.Green, m32.Blue
	Yellow, Cyan, Magenta            = m32.Yellow, m32.Cyan, m32.Magenta

	Clamp, Clamp01, Lerp, InvLerp, Remap = m32.Clamp, m32.Clamp01, m32.Lerp, m32.InvLerp, m32.Remap
	Deg2Rad, Rad2Deg                     = m32.Deg2Rad, m32.Rad2Deg
	Min, Max, Abs                        = m32.Min, m32.Max, m32.Abs

	Identity     = m32.Identity
	Translate    = m32.Translate
	Scale        = m32.Scale
	ScaleUniform = m32.ScaleUniform
	Rotate       = m32.Rotate
	RotateX      = m32.RotateX
	RotateY      = m32.RotateY
	RotateZ      = m32.RotateZ
	Perspective  = m32.Perspective
	Ortho        = m32.Ortho
	LookAt       = m32.LookAt
	QuatIdentity = m32.QuatIdentity
	QuatFromAxis = m32.QuatFromAxisAngle
)

// --- geometry -------------------------------------------------------------

type (
	Mesh   = geometry.Mesh
	Vertex = geometry.Vertex
)

var (
	NewMesh     = geometry.New
	Cube        = geometry.Cube
	Sphere      = geometry.Sphere
	Plane       = geometry.Plane
	Torus       = geometry.Torus
	Cylinder    = geometry.Cylinder
	Cone        = geometry.Cone
	Tetrahedron = geometry.Tetrahedron
	Icosahedron = geometry.Icosahedron
	LoadOBJ     = geometry.LoadOBJ
	LoadOBJFile = geometry.LoadOBJFile
)

// --- texture --------------------------------------------------------------

type Texture = texture.Texture

var (
	NewTexture  = texture.New
	FromImage   = texture.FromImage
	TexSolid    = texture.Solid
	TexChecker  = texture.Checker
	TexGradient = texture.Gradient
	TexNoise    = texture.Noise
	TexBrick    = texture.Brick
	TexFBM      = texture.FBM
)

// --- scene ----------------------------------------------------------------

type (
	Scene    = scene.Scene
	Node     = scene.Node
	Entity   = scene.Entity
	Material = scene.Material
	Light    = scene.Light
	Camera   = scene.Camera
	Fog      = scene.Fog
	Point    = scene.Point
)

var (
	NewScene    = scene.New
	NewNode     = scene.NewNode
	NewEntity   = scene.NewEntity
	NewMaterial = scene.NewMaterial
	NewCamera   = scene.NewCamera

	Ambient       = scene.Ambient
	Directional   = scene.Directional
	PointLight    = scene.PointLight
	Spot          = scene.Spot
	DefaultLights = scene.DefaultLights
)

// Shade modes.
const (
	ShadeFlat    = scene.ShadeFlat
	ShadeGouraud = scene.ShadeGouraud
	ShadePhong   = scene.ShadePhong
)

// Fog modes.
const (
	FogNone   = scene.FogNone
	FogLinear = scene.FogLinear
	FogExp    = scene.FogExp
	FogExp2   = scene.FogExp2
)

// --- render ---------------------------------------------------------------

type (
	Framebuffer = render.Framebuffer
	Renderer    = render.Renderer
)

var (
	NewFramebuffer = render.NewFramebuffer
	NewRenderer    = render.New
	SaveGIF        = render.SaveGIF
	EncodeGIF      = render.EncodeGIF
)

// --- anim -----------------------------------------------------------------

type (
	Tween          = anim.Tween
	Keyframe       = anim.Keyframe
	Vec3Track      = anim.Vec3Track
	Clip           = anim.Clip
	Animator       = anim.Animator
	GameLoop       = anim.GameLoop
	ParticleSystem = anim.ParticleSystem
	ParticleConfig = anim.ParticleConfig
	EaseFunc       = anim.EaseFunc
)

var (
	NewTween     = anim.NewTween
	NewClip      = anim.NewClip
	NewAnimator  = anim.NewAnimator
	NewGameLoop  = anim.NewGameLoop
	NewParticles = anim.NewParticleSystem

	EaseLinear     = anim.EaseLinear
	EaseInQuad     = anim.EaseInQuad
	EaseOutQuad    = anim.EaseOutQuad
	EaseInOutQuad  = anim.EaseInOutQuad
	EaseInCubic    = anim.EaseInCubic
	EaseOutCubic   = anim.EaseOutCubic
	EaseInOutCubic = anim.EaseInOutCubic
	EaseInSine     = anim.EaseInSine
	EaseOutSine    = anim.EaseOutSine
	EaseOutElastic = anim.EaseOutElastic
	EaseOutBounce  = anim.EaseOutBounce
)

// Render is a one-shot convenience: draw s through cam into fb.
func Render(s *Scene, fb *Framebuffer, cam *Camera) {
	NewRenderer(fb, cam).Render(s)
}
