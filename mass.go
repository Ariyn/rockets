package rockets

import (
	"math"
)

const gravity = 9.8

type MassObject struct {
	ID string

	IsKinematic bool
	Mass        float64
	// TODO: implement reverseMass equals to (1 / Mass)
	//reverseMass float64
	R float64

	PositionHistory []Vector3D
	Position        Vector3D
	Velocity        Vector3D
	Acceleration    Vector3D
	Force           Vector3D
	NextForce       Vector3D

	// TODO: implement these fields
	Angle      Vector3D
	AngularVel Vector3D
	AngularAcc Vector3D
	Torque     Vector3D
	NextTorque Vector3D
}

func (m *MassObject) Step(dt float64) {
	if m.IsKinematic {
		return
	}

	m.Force = m.NextForce
	m.Acceleration = m.Force.Mul(1 / m.Mass)
	m.Velocity = m.Velocity.Add(m.Acceleration.Mul(dt))
	m.Position = m.Position.Add(m.Velocity.Mul(dt))
	m.PositionHistory = append(m.PositionHistory, m.Position)
	m.NextForce = ZeroVector

	m.Torque = m.NextTorque
	m.AngularAcc = m.Torque.Mul(1 / m.Mass)
	m.AngularVel = m.AngularVel.Add(m.AngularAcc.Mul(dt))

	radZ := m.AngularVel.Z * (math.Pi / 180)
	if radZ != 0 {
		cr, sr := math.Cos(radZ), math.Sin(radZ)
		x := m.Angle.X*cr - m.Angle.Y*sr
		y := m.Angle.X*sr + m.Angle.Y*cr
		z := m.Angle.Z

		m.Angle = Vector3D{x, y, z}
	}

	radY := m.AngularVel.Y * (math.Pi / 180)
	if radY != 0 {
		cr, sr := math.Cos(radY), math.Sin(radY)
		x := m.Angle.X*cr + m.Angle.Z*sr
		y := m.Angle.Y
		z := -m.Angle.X*sr + m.Angle.Z*cr
		m.Angle = Vector3D{x, y, z}
	}

	radX := m.AngularVel.X * (math.Pi / 180)
	if radX != 0 {
		cr, sr := math.Cos(radX), math.Sin(radX)
		x := m.Angle.X
		y := m.Angle.Y*cr - m.Angle.Z*sr
		z := m.Angle.Y*sr + m.Angle.Z*cr
		m.Angle = Vector3D{x, y, z}
	}
	m.Torque = ZeroVector
}

func (m *MassObject) CalculateForce(other *MassObject) {
	if m.IsKinematic {
		return
	}

	p := other.Position.Sub(m.Position)
	r := p.Length()
	m.NextForce = m.NextForce.Add(p.Normalize().Mul(gravity * m.Mass * other.Mass / (r * r)))
}
