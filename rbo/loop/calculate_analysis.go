package main

func AnalyseData(intelligent, stupid [][]float64, factor float64) []Analysis {
	var analysis []Analysis
	for i := range intelligent[0] {
		a := CreateAnalysis(i,
			intelligent[0][i], intelligent[1][i], false,
			stupid[0][i], stupid[1][i], false,
			false, false, false, 0.0)
		analysis = append(analysis, a)
	}

	for i, a := range analysis {
		a.Intelligent.Stabil = (a.Intelligent.Mean > factor*a.Intelligent.StandardDeviation)
		a.Stupid.Stabil = (a.Stupid.Mean > factor*a.Stupid.StandardDeviation)
		analysis[i] = a
	}

	for i, a := range analysis {
		a.GoodMC = (a.Intelligent.Stabil == true && a.Stupid.Stabil == false)
		a.BadMC = (a.Intelligent.Stabil == false && a.Stupid.Stabil == true)
		analysis[i] = a
	}

	for i, a := range analysis {
		a.GoodMC = (a.Intelligent.Stabil == true && a.Stupid.Stabil == false)
		a.BadMC = (a.Intelligent.Stabil == false && a.Stupid.Stabil == true)
		analysis[i] = a
	}

	for i, a := range analysis {
		if a.Intelligent.Stabil == true && a.Stupid.Stabil == true {
			a.UseChange = true
			a.Change = a.Intelligent.Mean - a.Stupid.Mean
		}
		analysis[i] = a
	}

	return analysis
}
