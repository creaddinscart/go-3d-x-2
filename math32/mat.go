package math32

import "math"

// Mat4 is a 4x4 matrix stored in column-major order. Element (row r, col c)
// is m[c*4+r].
type Mat4 [16]float32

func (m *Mat4) at(r, c int) float32     { return m[c*4+r] }
func (m *Mat4) set(r, c int, v float32) { m[c*4+r] = v }

// At returns the element at row r, column c (0-indexed).
func (m Mat4) At(r, c int) float32 { return m[c*4+r] }

// Identity returns the 4x4 identity matrix.
func Identity() Mat4 {
	var m Mat4
	m[0], m[5], m[10], m[15] = 1, 1, 1, 1
	return m
}

// Translate builds a translation matrix.
func Translate(v Vec3) Mat4 {
	m := Identity()
	m.set(0, 3, v.X)
	m.set(1, 3, v.Y)
	m.set(2, 3, v.Z)
	return m
}

// Scale builds a non-uniform scale matrix.
func Scale(v Vec3) Mat4 {
	m := Identity()
	m.set(0, 0, v.X)
	m.set(1, 1, v.Y)
	m.set(2, 2, v.Z)
	return m
}

// ScaleUniform builds a uniform scale matrix.
func ScaleUniform(s float32) Mat4 { return Scale(Vec3{s, s, s}) }

// Rotate builds an axis-angle rotation matrix (right-handed, radians).
func Rotate(axis Vec3, angleRad float32) Mat4 {
	a := axis.Normalize()
	c := float32(math.Cos(float64(angleRad)))
	s := float32(math.Sin(float64(angleRad)))
	t := 1 - c
	x, y, z := a.X, a.Y, a.Z
	var m Mat4
	m.set(0, 0, t*x*x+c)
	m.set(0, 1, t*x*y-s*z)
	m.set(0, 2, t*x*z+s*y)
	m.set(1, 0, t*x*y+s*z)
	m.set(1, 1, t*y*y+c)
	m.set(1, 2, t*y*z-s*x)
	m.set(2, 0, t*x*z-s*y)
	m.set(2, 1, t*y*z+s*x)
	m.set(2, 2, t*z*z+c)
	m.set(3, 3, 1)
	return m
}

// RotateX rotates about the X axis.
func RotateX(a float32) Mat4 { return Rotate(UnitX, a) }

// RotateY rotates about the Y axis.
func RotateY(a float32) Mat4 { return Rotate(UnitY, a) }

// RotateZ rotates about the Z axis.
func RotateZ(a float32) Mat4 { return Rotate(UnitZ, a) }

// Mul multiplies two matrices (a * b).
func (a Mat4) Mul(b Mat4) Mat4 {
	var out Mat4
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			out.set(r, c, a.at(r, 0)*b.at(0, c)+a.at(r, 1)*b.at(1, c)+a.at(r, 2)*b.at(2, c)+a.at(r, 3)*b.at(3, c))
		}
	}
	return out
}

// MulVec4 transforms a homogeneous vector.
func (a Mat4) MulVec4(v Vec4) Vec4 {
	return Vec4{
		a.at(0, 0)*v.X + a.at(0, 1)*v.Y + a.at(0, 2)*v.Z + a.at(0, 3)*v.W,
		a.at(1, 0)*v.X + a.at(1, 1)*v.Y + a.at(1, 2)*v.Z + a.at(1, 3)*v.W,
		a.at(2, 0)*v.X + a.at(2, 1)*v.Y + a.at(2, 2)*v.Z + a.at(2, 3)*v.W,
		a.at(3, 0)*v.X + a.at(3, 1)*v.Y + a.at(3, 2)*v.Z + a.at(3, 3)*v.W,
	}
}

// TransformPoint applies the matrix as a position (w=1, includes translation).
func (a Mat4) TransformPoint(v Vec3) Vec3 { return a.MulVec4(v.ToVec4(1)).XYZ() }

// TransformVector applies the matrix as a direction (w=0, ignores translation).
func (a Mat4) TransformVector(v Vec3) Vec3 { return a.MulVec4(v.ToVec4(0)).XYZ() }

// Transpose returns the transpose.
func (a Mat4) Transpose() Mat4 {
	var out Mat4
	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			out.set(r, c, a.at(c, r))
		}
	}
	return out
}

// Inverse returns the general 4x4 inverse (identity fallback if singular),
// via the canonical cofactor expansion over column-major storage.
func (a Mat4) Inverse() Mat4 {
	a00, a01, a02, a03 := a[0], a[1], a[2], a[3]
	a10, a11, a12, a13 := a[4], a[5], a[6], a[7]
	a20, a21, a22, a23 := a[8], a[9], a[10], a[11]
	a30, a31, a32, a33 := a[12], a[13], a[14], a[15]

	b00 := a00*a11 - a01*a10
	b01 := a00*a12 - a02*a10
	b02 := a00*a13 - a03*a10
	b03 := a01*a12 - a02*a11
	b04 := a01*a13 - a03*a11
	b05 := a02*a13 - a03*a12
	b06 := a20*a31 - a21*a30
	b07 := a20*a32 - a22*a30
	b08 := a20*a33 - a23*a30
	b09 := a21*a32 - a22*a31
	b10 := a21*a33 - a23*a31
	b11 := a22*a33 - a23*a32

	det := b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06
	if det == 0 {
		return Identity()
	}
	inv := 1 / det
	return Mat4{
		(a11*b11 - a12*b10 + a13*b09) * inv,
		(-a01*b11 + a02*b10 - a03*b09) * inv,
		(a31*b05 - a32*b04 + a33*b03) * inv,
		(-a21*b05 + a22*b04 - a23*b03) * inv,
		(-a10*b11 + a12*b08 - a13*b07) * inv,
		(a00*b11 - a02*b08 + a03*b07) * inv,
		(-a30*b05 + a32*b02 - a33*b01) * inv,
		(a20*b05 - a22*b02 + a23*b01) * inv,
		(a10*b10 - a11*b08 + a13*b06) * inv,
		(-a00*b10 + a01*b08 - a03*b06) * inv,
		(a30*b04 - a31*b02 + a33*b00) * inv,
		(-a20*b04 + a21*b02 - a23*b00) * inv,
		(-a10*b09 + a11*b07 - a12*b06) * inv,
		(a00*b09 - a01*b07 + a02*b06) * inv,
		(-a30*b03 + a31*b01 - a32*b00) * inv,
		(a20*b03 - a21*b01 + a22*b00) * inv,
	}
}

// Perspective builds a right-handed perspective projection matrix.
// fovY is the vertical field of view in radians; near/far are positive.
func Perspective(fovY, aspect, near, far float32) Mat4 {
	f := float32(1) / float32(math.Tan(float64(fovY)/2))
	var m Mat4
	m.set(0, 0, f/aspect)
	m.set(1, 1, f)
	m.set(2, 2, (far+near)/(near-far))
	m.set(2, 3, (2*far*near)/(near-far))
	m.set(3, 2, -1)
	return m
}

// Ortho builds an orthographic projection matrix.
func Ortho(left, right, bottom, top, near, far float32) Mat4 {
	var m Mat4
	m.set(0, 0, 2/(right-left))
	m.set(1, 1, 2/(top-bottom))
	m.set(2, 2, -2/(far-near))
	m.set(0, 3, -(right+left)/(right-left))
	m.set(1, 3, -(top+bottom)/(top-bottom))
	m.set(2, 3, -(far+near)/(far-near))
	m.set(3, 3, 1)
	return m
}

// LookAt builds a right-handed view matrix.
func LookAt(eye, target, up Vec3) Mat4 {
	f := target.Sub(eye).Normalize()
	s := f.Cross(up.Normalize())
	u := s.Cross(f)
	var m Mat4
	m.set(0, 0, s.X)
	m.set(0, 1, s.Y)
	m.set(0, 2, s.Z)
	m.set(1, 0, u.X)
	m.set(1, 1, u.Y)
	m.set(1, 2, u.Z)
	m.set(2, 0, -f.X)
	m.set(2, 1, -f.Y)
	m.set(2, 2, -f.Z)
	m.set(0, 3, -s.Dot(eye))
	m.set(1, 3, -u.Dot(eye))
	m.set(2, 3, f.Dot(eye))
	m.set(3, 3, 1)
	return m
}
