package main

import "fmt"

func ConvertIROSAnalysisResults(dir, output string, analysis []Analysis) {
	var str []string

	for _, a := range analysis {
		if a.GoodMC == true {
			str = append(str, fmt.Sprintf("Good MC for Index %d with Coefficient %.3f\n", a.Index, a.Intelligent.Mean))
		}
		if a.BadMC == true {
			str = append(str, fmt.Sprintf("Bad MC for Index %d with Coefficient %.3f\n", a.Index, a.Stupid.Mean))
		}
		if a.UseChange == true {
			if a.Change > 0.0 {
				str = append(str, fmt.Sprintf("Make Index %d stiffer, because the difference between Good and Bad is %.3f\n", a.Index, a.Change))
			}
			if a.Change < 0.0 {
				str = append(str, fmt.Sprintf("Make Index %d more compliant, because the difference between Good and Bad is %.3f\n", a.Index, a.Change))
			}
		}
	}

	WriteStrings(fmt.Sprintf("%s/%s", dir, output), str)

}
