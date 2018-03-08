package main

import (
	"github.com/gonum/matrix/mat64"
	"github.com/westphae/quaternion"
)

type Pose struct {
	Position   P3D
	Quaternion quaternion.Quaternion
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

func PoseSub(src, target Pose) Pose {
	srcP := src.Position
	targetP := target.Position
	srcQ := src.Quaternion
	targetQ := target.Quaternion

	newTargetPosition := P3DSub(targetP, srcP)

	srcQInv := quaternion.Inv(srcQ)
	newTargetQuaternion := quaternion.Prod(srcQInv, targetQ) // relative rotation from src to target

	m := quaternion.RotMat(srcQ)

	var rotInv mat64.Dense
	var vRot mat64.Dense
	rot := mat64.NewDense(3, 3, []float64{m[0][0], m[0][1], m[0][2], m[1][0], m[1][1], m[1][2], m[2][0], m[2][1], m[2][2]})
	v := mat64.NewDense(3, 1, []float64{newTargetPosition.X, newTargetPosition.Y, newTargetPosition.Z})

	rotInv.Inverse(rot)
	vRot.Mul(&rotInv, v)

	return Pose{Position: P3D{X: vRot.At(0, 0), Y: vRot.At(1, 0), Z: vRot.At(2, 0)}, Quaternion: quaternion.Quaternion{X: newTargetQuaternion.X, Y: newTargetQuaternion.Y, Z: newTargetQuaternion.Z, W: newTargetQuaternion.W}}
}
