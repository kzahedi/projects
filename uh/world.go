package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	north = 0
	east  = 1
	south = 2
	west  = 3
)

type label struct {
	north int
	east  int
	south int
	west  int
}

func stringOfLabel(label int) string {
	switch label {
	case north:
		return "north"
	case east:
		return "east"
	case south:
		return "south"
	case west:
		return "west"
	}
	return "unknown"
}

func (l *label) String() string {
	s := ""
	s = fmt.Sprintf("north is labelled %s", stringOfLabel(l.north))
	s = fmt.Sprintf("%s, east is labelled %s", s, stringOfLabel(l.east))
	s = fmt.Sprintf("%s, south is labelled %s", s, stringOfLabel(l.south))
	s = fmt.Sprintf("%s, west is labelled %s", s, stringOfLabel(l.west))
	return s
}

type world struct {
	currentPosition int
	labelling       []label
}

func (w world) String() string {
	s := ""
	s = fmt.Sprintf("Agent is in state %d\n", w.currentPosition)
	for _, l := range w.labelling {
		s = fmt.Sprintf("%s%s\n", s, l.String())
	}
	return s
}

func createWorld(n int, t float64) world {
	w := world{currentPosition: 0, labelling: make([]label, n*n, n*n)}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < n*n; i++ {
		v := []int{0, 1, 2, 3}
		if rand.Float64() <= t {
			v = r.Perm(4)
		}
		w.labelling[i] = label{north: v[0], east: v[1], south: v[2], west: v[3]}
	}

	fmt.Println(w)
	return w
}
