package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kzahedi/goent/continuous"
	"github.com/kzahedi/goent/continuous/state"
	"github.com/kzahedi/utils"
)

func main() {

	prefix := flag.String("p", "ca", "prefix")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *help == true {
		flag.PrintDefaults()
		os.Exit(0)
	}

	aFile := fmt.Sprintf("%s_bizeps_trizeps_filtered.csv", *prefix)
	wFile := fmt.Sprintf("%s_signal.csv", *prefix)
	lFile := fmt.Sprintf("%s_labels.csv", *prefix)

	aRaw, _ := utils.ReadFloatCsv(aFile)
	wRaw, _ := utils.ReadFloatCsv(wFile)
	l, _ := utils.ReadCsv(lFile)

	marker := []string{"RHUMC", "RULN", "RELBW", "RSHO", "RRAD", "RHUMS", "RELB"}

	var indices []int

	for i, v := range l {
		for _, w := range marker {
			if w == v[0] {
				fmt.Println("Found ", w)
				indices = append(indices, 3*i+0)
				indices = append(indices, 3*i+1)
				indices = append(indices, 3*i+2)
			}
		}
	}

	wSel := utils.GetFloatColumns(wRaw, indices)

	wNormalised := continuous.Normalise(wSel)
	aNormalised := continuous.Normalise(aRaw)

	w2w1a1 := makeData(wNormalised, aNormalised)

	w2Indices := make([]int, len(wNormalised[0]), len(wNormalised[0]))
	w1Indices := make([]int, len(wNormalised[0]), len(wNormalised[0]))
	a1Indices := make([]int, len(aNormalised[0]), len(aNormalised[0]))

	for w := 0; w < len(wNormalised[0]); w++ {
		w2Indices[w] = w
		w1Indices[w] = len(wNormalised[0]) + w
	}
	for a := 0; a < len(aNormalised[0]); a++ {
		a1Indices[a] = 2*len(wNormalised[0]) + a
	}

	mcw := continuous.MorphologicalComputationW(w2w1a1, w2Indices, w1Indices, a1Indices, 30, true)
	mcw_pw := state.MorphologicalComputationW(w2w1a1, w2Indices, w1Indices, a1Indices, 30, true)

	fmt.Println(fmt.Sprintf("MC_W %f", mcw))

	output := fmt.Sprintf("%s_mcw.csv", *prefix)
	utils.WriteCsvFloatArray(output, mcw_pw, nil)
}

func makeData(w, a [][]float64) [][]float64 {
	nW := len(w[0])
	nA := len(a[0])
	rows := len(w) - 1
	cols := 2*nW + nA

	var r [][]float64

	for row := 0; row < rows; row++ {
		s := make([]float64, cols, cols)
		for wi := 0; wi < nW; wi++ {
			s[wi] = w[row+1][wi]
			s[nW+wi] = w[row][wi]
		}
		for ai := 0; ai < nA; ai++ {
			s[2*nW+ai] = a[row][ai]
		}
		r = append(r, s)
	}
	return r
}
