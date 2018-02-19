package main

import (
	"testing"
)

func TestP3DCopy(t *testing.T) {
	a := P3D{X: 1.0, Y: 2.0, Z: 3.0}
	b := P3DCopy(a)

	if a.X != b.X {
		t.Errorf("Copy of X failed. Should be %f, but is %f", a.X, b.X)
	}
	if a.Y != b.Y {
		t.Errorf("Copy of Y failed. Should be %f, but is %f", a.Y, b.Y)
	}
	if a.Z != b.Z {
		t.Errorf("Copy of Z failed. Should be %f, but is %f", a.Z, b.Z)
	}

	b.X = 10.0
	if a.X == b.X {
		t.Errorf("Copy of X failed. Both values are equal %f = %f", a.X, b.X)
	}
}

func TestP3DAdd(t *testing.T) {
	a := P3D{X: 1.0, Y: 2.0, Z: 3.0}
	b := P3D{X: 1.0, Y: 2.0, Z: 3.0}

	b.Add(a)

	if b.X != 2.0 {
		t.Errorf("Add failed, b.X should be 2.0 but it is %f", b.X)
	}
	if b.Y != 4.0 {
		t.Errorf("Add failed, b.Y should be 4.0 but it is %f", b.Y)
	}
	if b.Z != 6.0 {
		t.Errorf("Add failed, b.Z should be 6.0 but it is %f", b.Z)
	}

	if a.X != 1.0 {
		t.Errorf("Add failed, a.X should be 2.0 but it is %f", a.X)
	}
	if a.Y != 2.0 {
		t.Errorf("Add failed, a.Y should be 4.0 but it is %f", a.Y)
	}
	if a.Z != 3.0 {
		t.Errorf("Add failed, a.Z should be 6.0 but it is %f", a.Z)
	}
}

func TestP3DSub(t *testing.T) {
	a := P3D{X: 2.0, Y: 3.0, Z: 4.0}
	b := P3D{X: 1.0, Y: 2.0, Z: 3.0}

	c := P3DSub(a, b)

	if c.X != 1.0 {
		t.Errorf("Add failed, c.X should be 1.0 but it is %f", c.X)
	}
	if c.Y != 1.0 {
		t.Errorf("Add failed, c.Y should be 1.0 but it is %f", c.Y)
	}
	if c.Z != 1.0 {
		t.Errorf("Add failed, c.Z should be 1.0 but it is %f", c.Z)
	}
}

func TestPose(t *testing.T) {
	a := CreatePose(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)

	if a.Position.X != 1.0 {
		t.Errorf("position X should be 1.0 but it is %f", a.Position.X)
	}
	if a.Position.Y != 2.0 {
		t.Errorf("position Y should be 2.0 but it is %f", a.Position.Y)
	}
	if a.Position.Z != 3.0 {
		t.Errorf("position Z should be 3.0 but it is %f", a.Position.Z)
	}

	if a.Orientation.X != 4.0 {
		t.Errorf("position X should be 4.0 but it is %f", a.Orientation.X)
	}
	if a.Orientation.Y != 5.0 {
		t.Errorf("position Y should be 5.0 but it is %f", a.Orientation.Y)
	}
	if a.Orientation.Z != 6.0 {
		t.Errorf("position Z should be 6.0 but it is %f", a.Orientation.Z)
	}

}

func TestPoseCopy(t *testing.T) {
	a := CreatePose(1.0, 2.0, 3.0, 4.0, 5.0, 6.0)
	b := PoseCopy(a)

	if b.Position.X != 1.0 {
		t.Errorf("position X should be 1.0 but it is %f", b.Position.X)
	}
	if b.Position.Y != 2.0 {
		t.Errorf("position Y should be 2.0 but it is %f", b.Position.Y)
	}
	if b.Position.Z != 3.0 {
		t.Errorf("position Z should be 3.0 but it is %f", b.Position.Z)
	}

	if b.Orientation.X != 4.0 {
		t.Errorf("position X should be 4.0 but it is %f", b.Orientation.X)
	}
	if b.Orientation.Y != 5.0 {
		t.Errorf("position Y should be 5.0 but it is %f", b.Orientation.Y)
	}
	if b.Orientation.Z != 6.0 {
		t.Errorf("position Z should be 6.0 but it is %f", b.Orientation.Z)
	}

	b.Position.X = 2.0
	b.Position.Y = 3.0
	b.Position.Z = 4.0
	b.Orientation.X = 5.0
	b.Orientation.Y = 6.0
	b.Orientation.Z = 7.0

	if b.Position.X != 2.0 {
		t.Errorf("position X should be 2.0 but it is %f", b.Position.X)
	}
	if b.Position.Y != 3.0 {
		t.Errorf("position Y should be 3.0 but it is %f", b.Position.Y)
	}
	if b.Position.Z != 4.0 {
		t.Errorf("position Z should be 4.0 but it is %f", b.Position.Z)
	}

	if b.Orientation.X != 5.0 {
		t.Errorf("position X should be 5.0 but it is %f", b.Orientation.X)
	}
	if b.Orientation.Y != 6.0 {
		t.Errorf("position Y should be 6.0 but it is %f", b.Orientation.Y)
	}
	if b.Orientation.Z != 7.0 {
		t.Errorf("position Z should be 7.0 but it is %f", b.Orientation.Z)
	}

	if a.Position.X != 1.0 {
		t.Errorf("position X should be 1.0 but it is %f", a.Position.X)
	}
	if a.Position.Y != 2.0 {
		t.Errorf("position Y should be 2.0 but it is %f", a.Position.Y)
	}
	if a.Position.Z != 3.0 {
		t.Errorf("position Z should be 3.0 but it is %f", a.Position.Z)
	}

	if a.Orientation.X != 4.0 {
		t.Errorf("position X should be 4.0 but it is %f", a.Orientation.X)
	}
	if a.Orientation.Y != 5.0 {
		t.Errorf("position Y should be 5.0 but it is %f", a.Orientation.Y)
	}
	if a.Orientation.Z != 6.0 {
		t.Errorf("position Z should be 6.0 but it is %f", a.Orientation.Z)
	}

}
