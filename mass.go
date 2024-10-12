package rockets

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
	//Angle      Vector3D
	//AngularVel Vector3D
	//AngularAcc Vector3D
	//Torque     Vector3D
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
}

func (m *MassObject) CalculateForce(other *MassObject) {
	if m.IsKinematic {
		return
	}

	p := other.Position.Sub(m.Position)
	r := p.Length()
	m.NextForce = m.NextForce.Add(p.Normalize().Mul(gravity * m.Mass * other.Mass / (r * r)))
}
