package main

import (
	"fmt"

	"github.com/gonum/matrix/mat64"
	"github.com/westphae/quaternion"
)

type Result struct {
	MC_W           float64
	GraspDistance  float64
	Point          []float64
	ObjectType     int
	ObjectPosition int
	ClusteredByTSE bool
	Successful     bool
}

type Results map[string]Result

type P3D struct {
	X float64
	Y float64
	Z float64
}

type Pose struct {
	Position   P3D
	Quaternion quaternion.Quaternion
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

func QCopy(q quaternion.Quaternion) quaternion.Quaternion {
	return quaternion.Quaternion{X: q.X, Y: q.Y, Z: q.Z, W: q.W}
}

func PoseCopy(a Pose) Pose {
	return Pose{Position: P3DCopy(a.Position), Quaternion: QCopy(a.Quaternion)}
}

func CreatePose(x, y, z, qx, qy, qz, qw float64) Pose {
	return Pose{Position: P3D{X: x, Y: y, Z: z}, Quaternion: quaternion.Quaternion{X: qx, Y: qy, Z: qz, W: qw}}
}

func PoseSub(a, b Pose) Pose {
	aPos := a.Position
	bPos := b.Position
	cP := P3DSub(aPos, bPos)
	aQ := quaternion.Inv(a.Quaternion)
	bQ := b.Quaternion
	cQ := quaternion.Prod(aQ, bQ)
	m := quaternion.RotMat(aQ)

	var rotInv mat64.Dense
	var vRot mat64.Dense
	rot := mat64.NewDense(3, 3, []float64{m[0][0], m[0][1], m[0][2], m[1][0], m[1][1], m[1][2], m[2][0], m[2][1], m[2][2]})
	v := mat64.NewDense(3, 1, []float64{cP.X, cP.Y, cP.Z})

	rotInv.Inverse(rot)
	vRot.Mul(&rotInv, v)

	return Pose{Position: P3D{X: vRot.At(0, 0), Y: vRot.At(1, 0), Z: vRot.At(2, 0)}, Quaternion: quaternion.Quaternion{X: cQ.X, Y: cQ.Y, Z: cQ.Z, W: cQ.W}}
}

func PrintResults(r map[string]Result) {
	for key, value := range r {
		fmt.Println(fmt.Sprintf("%s: MC_W: %f, Grasp Distance: %f, Point: (%f,%f), Object Type: %d, Object Position %d", key, value.MC_W, value.GraspDistance, value.Point[0], value.Point[1], value.ObjectType, value.ObjectPosition))
	}
}
