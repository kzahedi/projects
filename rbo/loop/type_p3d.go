package main

type P3D struct {
	X float64
	Y float64
	Z float64
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
