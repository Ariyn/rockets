package rockets

type MassObject struct {
	ID   string
	Mass float64
	// TODO: implement reverseMass equals to (1 / Mass)
	//reverseMass float64
	R float64

	PositionHistory []Vector3D
	Position        Vector3D
	Velocity        Vector3D
	Acceleration    Vector3D
	Force           Vector3D

	// TODO: implement these fields
	//Angle      Vector3D
	//AngularVel Vector3D
	//AngularAcc Vector3D
	//Torque     Vector3D
}

func (m *MassObject) Step(dt float64) {
	m.PositionHistory = append(m.PositionHistory, m.Position)
	m.Position = m.Position.Add(m.Velocity.Mul(dt))
	m.Velocity = m.Velocity.Add(m.Acceleration.Mul(dt))
	m.Acceleration = m.Force.Mul(1 / m.Mass)
}
