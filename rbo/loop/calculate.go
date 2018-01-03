package main

// CalculateGlobalVelocities This function calculates the velocities and angular
// velocities in the world coordinate system, which means that
// velocity(t)         = position(t) - position(t-1) and
// angular velocity(t) = orientation(t) - orientation(t-1)
// velocity(0)         = 0
// angular velocity(0) = 0
func CalculateGlobalVelocities(data Data) Data {
	for trajectoryIndex := 0; trajectoryIndex < data.NrOfTrajectories; trajectoryIndex++ {
		data.Trajectories[trajectoryIndex].GlobalVelocity = make([]Pose, data.NrOfDataPoints, data.NrOfDataPoints)
		for frameIndex := 1; frameIndex < data.NrOfDataPoints; frameIndex++ {
			diff :=
				PoseSub(
					data.Trajectories[trajectoryIndex].Frame[frameIndex],
					data.Trajectories[trajectoryIndex].Frame[frameIndex-1])
			data.Trajectories[trajectoryIndex].GlobalVelocity[frameIndex] = diff
		}
	}
	return data
}
