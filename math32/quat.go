package math32

import "math"

// Quaternion represents a rotation.
type Quaternion struct{ X, Y, Z, W float32 }

// QuatIdentity returns the identity quaternion.
func QuatIdentity() Quaternion { return Quaternion{0, 0, 0, 1} }

// QuatFromAxisAngle builds a quaternion from an axis and angle in radians.
func QuatFromAxisAngle(axis Vec3, angleRad float32) Quaternion {
	a := axis.Normalize()
	h := angleRad / 2
	s := float32(math.Sin(float64(h)))
	return Quaternion{a.X * s, a.Y * s, a.Z * s, float32(math.Cos(float64(h)))}
}

// Normalize returns a unit quaternion.
func (q Quaternion) Normalize() Quaternion {
	l := float32(math.Sqrt(float64(q.X*q.X + q.Y*q.Y + q.Z*q.Z + q.W*q.W)))
	if l == 0 {
		return q
	}
	return Quaternion{q.X / l, q.Y / l, q.Z / l, q.W / l}
}

// Mul composes two rotations.
func (q Quaternion) Mul(o Quaternion) Quaternion {
	return Quaternion{
		q.W*o.X + q.X*o.W + q.Y*o.Z - q.Z*o.Y,
		q.W*o.Y - q.X*o.Z + q.Y*o.W + q.Z*o.X,
		q.W*o.Z + q.X*o.Y - q.Y*o.X + q.Z*o.W,
		q.W*o.W - q.X*o.X - q.Y*o.Y - q.Z*o.Z,
	}
}

// Mat4 converts the quaternion to a rotation matrix.
func (q Quaternion) Mat4() Mat4 {
	q = q.Normalize()
	xx, yy, zz := q.X*q.X, q.Y*q.Y, q.Z*q.Z
	xy, xz, yz := q.X*q.Y, q.X*q.Z, q.Y*q.Z
	wx, wy, wz := q.W*q.X, q.W*q.Y, q.W*q.Z
	var m Mat4
	m.set(0, 0, 1-2*(yy+zz))
	m.set(0, 1, 2*(xy-wz))
	m.set(0, 2, 2*(xz+wy))
	m.set(1, 0, 2*(xy+wz))
	m.set(1, 1, 1-2*(xx+zz))
	m.set(1, 2, 2*(yz-wx))
	m.set(2, 0, 2*(xz-wy))
	m.set(2, 1, 2*(yz+wx))
	m.set(2, 2, 1-2*(xx+yy))
	m.set(3, 3, 1)
	return m
}

// RotateVec rotates a vector by the quaternion.
func (q Quaternion) RotateVec(v Vec3) Vec3 {
	return q.Mat4().TransformVector(v)
}

// Slerp spherical-linearly interpolates between two quaternions.
func (q Quaternion) Slerp(o Quaternion, t float32) Quaternion {
	dot := q.X*o.X + q.Y*o.Y + q.Z*o.Z + q.W*o.W
	if dot < 0 {
		dot = -dot
		o = Quaternion{-o.X, -o.Y, -o.Z, -o.W}
	}
	if dot > 0.9995 {
		return Quaternion{
			Lerp(q.X, o.X, t), Lerp(q.Y, o.Y, t),
			Lerp(q.Z, o.Z, t), Lerp(q.W, o.W, t),
		}.Normalize()
	}
	theta := float32(math.Acos(float64(Clamp01(dot))))
	sinT := float32(math.Sin(float64(theta)))
	w1 := float32(math.Sin(float64((1-t)*theta))) / sinT
	w2 := float32(math.Sin(float64(t*theta))) / sinT
	return Quaternion{
		q.X*w1 + o.X*w2, q.Y*w1 + o.Y*w2,
		q.Z*w1 + o.Z*w2, q.W*w1 + o.W*w2,
	}
}
