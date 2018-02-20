package main

import "fmt"

type Result struct {
	MC_W           float64
	GraspDistance  float64
	Point          []float64
	ObjectType     int
	ObjectPosition int
	ClusteredByTSE bool
}

type Results map[string]Result

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

func P3DCopy(a P3D) P3D {
	return P3D{X: a.X, Y: a.Y, Z: a.Z}
}

func PoseCopy(a Pose) Pose {
	return Pose{Position: P3DCopy(a.Position), Orientation: P3DCopy(a.Orientation)}
}

func PoseSub(a, b Pose) Pose {
	return Pose{Position: P3DSub(a.Position, b.Position), Orientation: P3DSub(a.Orientation, b.Orientation)}
}

func CreatePose(x, y, z, alpha, beta, gamma float64) Pose {
	return Pose{Position: P3D{X: x, Y: y, Z: z}, Orientation: P3D{X: alpha, Y: beta, Z: gamma}}
}

func PrintResults(r map[string]Result) {
	for key, value := range r {
		fmt.Println(fmt.Sprintf("%s: MC_W: %f, Grasp Distance: %f, Point: (%f,%f), Object Type: %d, Object Position %d", key, value.MC_W, value.GraspDistance, value.Point[0], value.Point[1], value.ObjectType, value.ObjectPosition))
	}
}
