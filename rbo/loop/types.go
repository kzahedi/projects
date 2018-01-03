package main

type P3D struct {
	X float64
	Y float64
	Z float64
}

type Pose struct {
	Position    P3D
	Orientation P3D
}

type Trajectory struct {
	Frame          []Pose
	GlobalVelocity []Pose // pose(t) - pose(t-1)
	LocalVelocity  []Pose // pose(t) - pose(t-1) in local coordinate frame of preceding frame
}

type Data struct {
	Trajectories     []Trajectory
	NrOfDataPoints   int
	NrOfTrajectories int
}

func (p *P3D) Add(q P3D) {
	p.X += q.X
	p.Y += q.Y
	p.Z += q.Z
}

func P3DSub(a, b P3D) P3D {
	return P3D{X: a.X - b.X, Y: a.Y - b.Y, Z: a.Z - b.Z}
}

func PoseSub(a, b Pose) Pose {
	return Pose{Position: P3DSub(a.Position, b.Position), Orientation: P3DSub(a.Orientation, b.Orientation)}
}
