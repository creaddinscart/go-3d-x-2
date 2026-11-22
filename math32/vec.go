package math32

import "math"

// Vec2 is a 2-component float32 vector.
type Vec2 struct{ X, Y float32 }

// Vec3 is a 3-component float32 vector.
type Vec3 struct{ X, Y, Z float32 }

// Vec4 is a 4-component float32 vector / homogeneous coordinate.
type Vec4 struct{ X, Y, Z, W float32 }

// V2 builds a Vec2.
func V2(x, y float32) Vec2 { return Vec2{x, y} }

// V3 builds a Vec3.
func V3(x, y, z float32) Vec3 { return Vec3{x, y, z} }

// V4 builds a Vec4.
func V4(x, y, z, w float32) Vec4 { return Vec4{x, y, z, w} }

// Axis-aligned unit vectors.
var (
	Zero3 = Vec3{}
	UnitX = Vec3{1, 0, 0}
	UnitY = Vec3{0, 1, 0}
	UnitZ = Vec3{0, 0, 1}
	One3  = Vec3{1, 1, 1}
)

// Vec3 methods -------------------------------------------------------------

func (a Vec3) Add(b Vec3) Vec3     { return Vec3{a.X + b.X, a.Y + b.Y, a.Z + b.Z} }
func (a Vec3) Sub(b Vec3) Vec3     { return Vec3{a.X - b.X, a.Y - b.Y, a.Z - b.Z} }
func (a Vec3) Mul(s float32) Vec3  { return Vec3{a.X * s, a.Y * s, a.Z * s} }
func (a Vec3) Div(s float32) Vec3  { return Vec3{a.X / s, a.Y / s, a.Z / s} }
func (a Vec3) Neg() Vec3           { return Vec3{-a.X, -a.Y, -a.Z} }
func (a Vec3) Dot(b Vec3) float32  { return a.X*b.X + a.Y*b.Y + a.Z*b.Z }
func (a Vec3) LengthSq() float32   { return a.Dot(a) }
func (a Vec3) Length() float32     { return float32(math.Sqrt(float64(a.Dot(a)))) }
func (a Vec3) SetY(y float32) Vec3 { return Vec3{a.X, y, a.Z} }

// Distance returns the Euclidean distance between two points.
func (a Vec3) Distance(b Vec3) float32 {
	return float32(math.Sqrt(float64(a.Sub(b).LengthSq())))
}

// Cross returns the vector perpendicular to a and b.
func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

// Normalize returns a unit-length copy (zero vector stays zero).
func (a Vec3) Normalize() Vec3 {
	l := a.Length()
	if l == 0 {
		return a
	}
	return a.Div(l)
}

// Lerp linearly interpolates between a and b.
func (a Vec3) Lerp(b Vec3, t float32) Vec3 {
	return Vec3{Lerp(a.X, b.X, t), Lerp(a.Y, b.Y, t), Lerp(a.Z, b.Z, t)}
}

// ToVec4 extends the vector with a w component.
func (a Vec3) ToVec4(w float32) Vec4 { return Vec4{a.X, a.Y, a.Z, w} }

// XY drops the Z component.
func (a Vec3) XY() Vec2 { return Vec2{a.X, a.Y} }

// Vec2 methods -------------------------------------------------------------

func (a Vec2) Add(b Vec2) Vec2    { return Vec2{a.X + b.X, a.Y + b.Y} }
func (a Vec2) Sub(b Vec2) Vec2    { return Vec2{a.X - b.X, a.Y - b.Y} }
func (a Vec2) Mul(s float32) Vec2 { return Vec2{a.X * s, a.Y * s} }
func (a Vec2) Dot(b Vec2) float32 { return a.X*b.X + a.Y*b.Y }
func (a Vec2) Length() float32    { return float32(math.Sqrt(float64(a.X*a.X + a.Y*a.Y))) }
func (a Vec2) Lerp(b Vec2, t float32) Vec2 {
	return Vec2{Lerp(a.X, b.X, t), Lerp(a.Y, b.Y, t)}
}
func (a Vec2) ToVec3(z float32) Vec3 { return Vec3{a.X, a.Y, z} }

// Vec4 methods -------------------------------------------------------------

func (a Vec4) Add(b Vec4) Vec4    { return Vec4{a.X + b.X, a.Y + b.Y, a.Z + b.Z, a.W + b.W} }
func (a Vec4) Mul(s float32) Vec4 { return Vec4{a.X * s, a.Y * s, a.Z * s, a.W * s} }

// XYZ returns the perspective divide (x/w, y/w, z/w).
func (a Vec4) XYZ() Vec3 {
	if a.W == 0 {
		return Vec3{a.X, a.Y, a.Z}
	}
	return Vec3{a.X / a.W, a.Y / a.W, a.Z / a.W}
}
