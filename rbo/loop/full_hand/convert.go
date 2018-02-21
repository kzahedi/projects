package main

import "math"

func convertX(data Data, trajectoryIndex, frameIndex int) float64 {
	previous := data.Trajectories[trajectoryIndex].Frame[frameIndex-1].Orientation.X
	current := data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.X
	if previous < 0.0 && current > 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + current - math.Pi
	}
	if previous > 0.0 && current < 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + math.Pi + current
	}
	return current
}

func convertY(data Data, trajectoryIndex, frameIndex int) float64 {
	previous := data.Trajectories[trajectoryIndex].Frame[frameIndex-1].Orientation.Y
	current := data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Y
	if previous < 0.0 && current > 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + current - math.Pi
	}
	if previous > 0.0 && current < 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + math.Pi + current
	}
	return current
}

func convertZ(data Data, trajectoryIndex, frameIndex int) float64 {
	previous := data.Trajectories[trajectoryIndex].Frame[frameIndex-1].Orientation.Z
	current := data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Z
	if previous < 0.0 && current > 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + current - math.Pi
	}
	if previous > 0.0 && current < 0.0 && math.Abs(previous) > 0.1 && math.Abs(current) > 0.1 {
		current = previous + math.Pi + current
	}
	return current
}

// make angles continuous without jump from +pi to -pi and vice versa
func ConvertAngles(data Data) Data {
	for trajectoryIndex := 0; trajectoryIndex < data.NrOfTrajectories; trajectoryIndex++ {
		for frameIndex := 1; frameIndex < data.NrOfDataPoints; frameIndex++ {
			data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.X = convertX(data, trajectoryIndex, frameIndex)
			data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Y = convertY(data, trajectoryIndex, frameIndex)
			data.Trajectories[trajectoryIndex].Frame[frameIndex].Orientation.Z = convertZ(data, trajectoryIndex, frameIndex)
		}
	}
	return data
}
