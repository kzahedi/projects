package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type AnalysisValue struct {
	Mean              float64
	StandardDeviation float64
	Stabil            bool
}

type Analysis struct {
	Index       int
	Intelligent AnalysisValue
	Stupid      AnalysisValue
	GoodMC      bool
	BadMC       bool
	UseChange   bool
	Change      float64
}

func (a Analysis) String() string {
	return fmt.Sprintf("Index %d:\n  Mean (I)   %.3f\n  STD (I)    %.3f\n  Stabil (I) %t\n  Mean (S)   %.3f\n  STD (S)    %.3f\n  Stabil (S) %t\n  GoodMC %t\n  BadMC %t\n  Use Change %t\n  Change %.3f", a.Index, a.Intelligent.Mean, a.Intelligent.StandardDeviation, a.Intelligent.Stabil, a.Stupid.Mean, a.Stupid.StandardDeviation, a.Stupid.Stabil, a.GoodMC, a.BadMC, a.UseChange, a.Change)
}

func CreateAnalysis(index int, imean, istd float64, istabil bool, smean, sstd float64, sstabil bool, goodMC, badMC, useChange bool, change float64) Analysis {
	intelligent := AnalysisValue{Mean: imean, StandardDeviation: istd, Stabil: istabil}
	stupid := AnalysisValue{Mean: smean, StandardDeviation: sstd, Stabil: sstabil}
	return Analysis{Index: index, Intelligent: intelligent, Stupid: stupid,
		GoodMC: goodMC, BadMC: badMC, UseChange: useChange, Change: change}
}

func WriteAnalysis(dir, filename string, analysis []Analysis) {
	file, err := os.Create(fmt.Sprintf("%s/%s", dir, filename))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(file)

	s := "# Index, Intelligent.Mean, Intelligent.StandardDeviation, Intelligent.Stabil, Stupid.Mean, Stupid.StandardDeviation, Stupid.Stabil, Good MC, Bad MC, Use Change, Change"

	w.WriteString(s)

	for _, a := range analysis {
		s = fmt.Sprintf("\n%d,%f,%f,%t,%f,%f,%t,%t,%t,%t,%f", a.Index,
			a.Intelligent.Mean, a.Intelligent.StandardDeviation, a.Intelligent.Stabil,
			a.Stupid.Mean, a.Stupid.StandardDeviation, a.Stupid.Stabil,
			a.GoodMC, a.BadMC, a.UseChange, a.Change)
		w.WriteString(s)
		w.Flush()
	}
}
