package main

import "fmt"

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

func printFrames(t Trajectory) {
	for _, f := range t.Frame {
		fmt.Println(
			f.Position.X, " ", f.Position.Y, " ", f.Position.Z,
			f.Quaternion.X, " ", f.Quaternion.Y, " ", f.Quaternion.Z, " ", f.Quaternion.W)
	}
}
