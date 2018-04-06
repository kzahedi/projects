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
	var transformedPosition mat64.Dense
	srcP := src.Position
	srcQ := src.Quaternion
	srcRot := quaternion.RotMat(srcQ)
	srcTransformation := mat64.NewDense(4, 4, nil)
	invSrcTransformation := mat64.NewDense(4, 4, nil)

	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			srcTransformation.Set(i, j, srcRot[i][j])
		}
	}

	srcTransformation.Set(0, 3, srcP.X)
	srcTransformation.Set(1, 3, srcP.Y)
	srcTransformation.Set(2, 3, srcP.Z)
	srcTransformation.Set(3, 3, 1.0)

	invSrcTransformation.Inverse(srcTransformation)

	targetP := target.Position
	targetQ := target.Quaternion

	targetPositionVector := mat64.NewDense(4, 1, []float64{targetP.X, targetP.Y, targetP.Z, 0.0})

	srcQInv := quaternion.Inv(srcQ)
	newTargetQuaternion := quaternion.Prod(srcQInv, targetQ) // relative rotation from src to target

	transformedPosition.Mul(invSrcTransformation, targetPositionVector)
	// fmt.Println(transformedPosition)

	return Pose{Position: P3D{X: transformedPosition.At(0, 0), Y: transformedPosition.At(1, 0), Z: transformedPosition.At(2, 0)}, Quaternion: quaternion.Quaternion{X: newTargetQuaternion.X, Y: newTargetQuaternion.Y, Z: newTargetQuaternion.Z, W: newTargetQuaternion.W}}
}
