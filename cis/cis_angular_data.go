package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kzahedi/goent/continuous"
	cstate "github.com/kzahedi/goent/continuous/state"
	"github.com/kzahedi/goent/dh"
	"github.com/kzahedi/goent/discrete"
	dstate "github.com/kzahedi/goent/discrete/state"
	"github.com/kzahedi/utils"
)

func main() {

	prefix := flag.String("p", "ca", "prefix")
	discrete := flag.Bool("d", false, "discrete")
	bins := flag.Int("b", 100, "bins")
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

	labelAngle := []string{"RShoulderAnteRetroHDCalProjAngle", "RElbowFlexExtAngle"}

	var indices []int

	for i, v := range l {
		for _, w := range labelAngle {
			if w == v[0] {
				fmt.Println("Found ", w, " at index ", i)
				indices = append(indices, i)
			}
		}
	}

	wSel := utils.GetFloatColumns(wRaw, indices)

	if *discrete == true {
		computeDiscrete(wSel, aRaw, *bins, *prefix)
	} else {
		computeContinuous(wSel, aRaw, *prefix)
	}
}

func computeDiscrete(w, a [][]float64, bins int, prefix string) {

	wmin, wmax := dh.GetMinMax(w)
	amin, amax := dh.GetMinMax(a)

	b := make([]int, len(w[0]), len(w[0]))
	for i := range b {
		b[i] = bins
	}

	wD := dh.Discretise(w, b, wmin, wmax)
	aD := dh.Discretise(a, b, amin, amax)

	W := dh.MakeUnivariateRelabelled(wD, b)
	A := dh.MakeUnivariateRelabelled(aD, b)

	w2w1a1 := makeDataInt(W, A)

	pw2w1a1 := discrete.Emperical3D(w2w1a1)

	mcw := discrete.MorphologicalComputationW(pw2w1a1)
	mcw_pw := dstate.MorphologicalComputationW(w2w1a1)

	fmt.Println(fmt.Sprintf("MC_W %f", mcw))

	output := fmt.Sprintf("%s_mcw_discrete.csv", prefix)
	utils.WriteCsvFloatArray(output, mcw_pw, nil)
}

func computeContinuous(w, a [][]float64, prefix string) {
	wNormalised := continuous.Normalise(w)
	aNormalised := continuous.Normalise(a)

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
	mcw_pw := cstate.MorphologicalComputationW(w2w1a1, w2Indices, w1Indices, a1Indices, 30, true)

	fmt.Println(fmt.Sprintf("MC_W %f", mcw))

	output := fmt.Sprintf("%s_mcw_cont.csv", prefix)
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

func makeDataInt(w, a []int) (r [][]int) {
	for row := 0; row < len(w)-1; row++ {
		r = append(r, []int{w[row+1], w[row], a[row]})
	}
	return r
}
