package rockets

import "math"

var ZeroVector = Vector3D{0, 0, 0}

type Vector3D struct {
	X, Y, Z float64
}

func (v Vector3D) Normalize() Vector3D {
	return v.Mul(1 / v.Length())
}

func (v Vector3D) Sub(v2 Vector3D) Vector3D {
	return Vector3D{v.X - v2.X, v.Y - v2.Y, v.Z - v2.Z}
}

func (v Vector3D) Mul(scalar float64) Vector3D {
	return Vector3D{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vector3D) Add(v2 Vector3D) Vector3D {
	return Vector3D{v.X + v2.X, v.Y + v2.Y, v.Z + v2.Z}
}

func (v Vector3D) Dot(v2 Vector3D) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v Vector3D) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector3D) Angle(v2 Vector3D) float64 {
	return math.Cos(v.Dot(v2) / (v.Length() * v2.Length()))
}

func (v Vector3D) Distance(v2 Vector3D) float64 {
	return v.Sub(v2).Length()
}
