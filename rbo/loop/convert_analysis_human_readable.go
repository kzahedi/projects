package main

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// IROS
////////////////////////////////////////////////////////////////////////////////

func ConvertIROSAnalysisResults(dir, output string, analysis []Analysis) {
	var str []string

	thumb := [][]int{{29, 30}, {28, 29}, {27, 28}, {26, 27}, {25, 30}}
	palm := [][]int{{24, 25}, {23, 24}, {22, 23}, {21, 22}, {0, 21}}
	pinky := [][]int{{19, 20}, {18, 19}, {17, 18}, {16, 17}, {0, 16}}
	ring := [][]int{{14, 15}, {13, 14}, {12, 13}, {11, 12}, {0, 11}}
	middle := [][]int{{9, 10}, {8, 9}, {7, 8}, {6, 7}, {0, 6}}
	index := [][]int{{5, 6}, {4, 5}, {3, 4}, {2, 3}, {0, 1}}

	indexStrings := getAnalysisStringsIROS("Index Finger", index, analysis)
	middleStrings := getAnalysisStringsIROS("Middle Finger", middle, analysis)
	ringStrings := getAnalysisStringsIROS("Ring Finger", ring, analysis)
	pinkyStrings := getAnalysisStringsIROS("Pinky Finger", pinky, analysis)
	palmStrings := getAnalysisStringsIROS("Palm", palm, analysis)
	thumbStrings := getAnalysisStringsIROS("Thumb", thumb, analysis)

	str = append(str, "############################################################\n")
	str = append(str, "##### Index Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range indexStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Middle Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range middleStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Ring Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range ringStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Pinky Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range pinkyStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Palm Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range palmStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Thumb Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range thumbStrings {
		str = append(str, v)
	}

	WriteStrings(fmt.Sprintf("%s/%s", dir, output), str)
}

func getAnalysisStringsIROS(label string, indices [][]int, analysis []Analysis) []string {
	var r []string

	for _, v := range indices {
		src := v[0]
		dst := v[1]

		xSrcIndex := src*3 + 0
		ySrcIndex := src*3 + 1
		zSrcIndex := src*3 + 2

		xDstIndex := dst*3 + 0
		yDstIndex := dst*3 + 1
		zDstIndex := dst*3 + 2

		xxIndex := xSrcIndex*93 + xDstIndex
		xyIndex := xSrcIndex*93 + yDstIndex
		xzIndex := xSrcIndex*93 + zDstIndex

		yxIndex := ySrcIndex*93 + xDstIndex
		yyIndex := ySrcIndex*93 + yDstIndex
		yzIndex := ySrcIndex*93 + zDstIndex

		zxIndex := zSrcIndex*93 + xDstIndex
		zyIndex := zSrcIndex*93 + yDstIndex
		zzIndex := zSrcIndex*93 + zDstIndex

		for _, a := range analysis {
			r = appendGoodMCString(xxIndex, label, src, dst, "x vs. x", a, r)
			r = appendGoodMCString(xyIndex, label, src, dst, "x vs. y", a, r)
			r = appendGoodMCString(xzIndex, label, src, dst, "x vs. z", a, r)

			r = appendGoodMCString(yxIndex, label, src, dst, "y vs. x", a, r)
			r = appendGoodMCString(yyIndex, label, src, dst, "y vs. y", a, r)
			r = appendGoodMCString(yzIndex, label, src, dst, "y vs. z", a, r)

			r = appendGoodMCString(zxIndex, label, src, dst, "z vs. x", a, r)
			r = appendGoodMCString(zyIndex, label, src, dst, "z vs. y", a, r)
			r = appendGoodMCString(zzIndex, label, src, dst, "z vs. z", a, r)
		}
	}

	for _, v := range indices {
		src := v[0]
		dst := v[1]

		xSrcIndex := src*3 + 0
		ySrcIndex := src*3 + 1
		zSrcIndex := src*3 + 2

		xDstIndex := dst*3 + 0
		yDstIndex := dst*3 + 1
		zDstIndex := dst*3 + 2

		xxIndex := xSrcIndex*93 + xDstIndex
		xyIndex := xSrcIndex*93 + yDstIndex
		xzIndex := xSrcIndex*93 + zDstIndex

		yxIndex := ySrcIndex*93 + xDstIndex
		yyIndex := ySrcIndex*93 + yDstIndex
		yzIndex := ySrcIndex*93 + zDstIndex

		zxIndex := zSrcIndex*93 + xDstIndex
		zyIndex := zSrcIndex*93 + yDstIndex
		zzIndex := zSrcIndex*93 + zDstIndex

		for _, a := range analysis {
			r = appendBadMCString(xxIndex, label, src, dst, "x vs. x", a, r)
			r = appendBadMCString(xyIndex, label, src, dst, "x vs. y", a, r)
			r = appendBadMCString(xzIndex, label, src, dst, "x vs. z", a, r)

			r = appendBadMCString(yxIndex, label, src, dst, "y vs. x", a, r)
			r = appendBadMCString(yyIndex, label, src, dst, "y vs. y", a, r)
			r = appendBadMCString(yzIndex, label, src, dst, "y vs. z", a, r)

			r = appendBadMCString(zxIndex, label, src, dst, "z vs. x", a, r)
			r = appendBadMCString(zyIndex, label, src, dst, "z vs. y", a, r)
			r = appendBadMCString(zzIndex, label, src, dst, "z vs. z", a, r)
		}
	}

	for _, v := range indices {
		src := v[0]
		dst := v[1]

		xSrcIndex := src*3 + 0
		ySrcIndex := src*3 + 1
		zSrcIndex := src*3 + 2

		xDstIndex := dst*3 + 0
		yDstIndex := dst*3 + 1
		zDstIndex := dst*3 + 2

		xxIndex := xSrcIndex*93 + xDstIndex
		xyIndex := xSrcIndex*93 + yDstIndex
		xzIndex := xSrcIndex*93 + zDstIndex

		yxIndex := ySrcIndex*93 + xDstIndex
		yyIndex := ySrcIndex*93 + yDstIndex
		yzIndex := ySrcIndex*93 + zDstIndex

		zxIndex := zSrcIndex*93 + xDstIndex
		zyIndex := zSrcIndex*93 + yDstIndex
		zzIndex := zSrcIndex*93 + zDstIndex

		for _, a := range analysis {
			r = appendChangeString(xxIndex, label, src, dst, "x vs. x", a, r)
			r = appendChangeString(xyIndex, label, src, dst, "x vs. y", a, r)
			r = appendChangeString(xzIndex, label, src, dst, "x vs. z", a, r)

			r = appendChangeString(yxIndex, label, src, dst, "y vs. x", a, r)
			r = appendChangeString(yyIndex, label, src, dst, "y vs. y", a, r)
			r = appendChangeString(yzIndex, label, src, dst, "y vs. z", a, r)

			r = appendChangeString(zxIndex, label, src, dst, "z vs. x", a, r)
			r = appendChangeString(zyIndex, label, src, dst, "z vs. y", a, r)
			r = appendChangeString(zzIndex, label, src, dst, "z vs. z", a, r)
		}
	}
	return r
}

func appendGoodMCString(index int, label string, src, dst int, vs string, a Analysis, r []string) []string {
	if a.Index == index && a.GoodMC == true {
		return append(r, fmt.Sprintf("%s Frame %d vs. Frame %d (%s): Good MC with value %.3f\n", label, src, dst, vs, a.Intelligent.Mean))
	}
	return r
}

func appendBadMCString(index int, label string, src, dst int, vs string, a Analysis, r []string) []string {
	if a.Index == index && a.BadMC == true {
		return append(r, fmt.Sprintf("%s Frame %d vs. Frame %d (%s): Bad MC with value %.3f\n", label, src, dst, vs, a.Intelligent.Mean))
	}
	return r
}

func appendChangeString(index int, label string, src, dst int, vs string, a Analysis, r []string) []string {
	if a.Index == index && a.UseChange == true {
		rec := "increase stiffness"
		if a.Change < 0.0 {
			rec = "increase compliance"
		}
		return append(r, fmt.Sprintf("%s Frame %d vs. Frame %d (%s): Stable values where detected for Intelligent (mean %.3f, std %.3f) Stupid (mean %.3f, std %.3f). Difference (mean values) is %.3f. Recommendation is %s\n", label, src, dst, vs, a.Intelligent.Mean, a.Intelligent.StandardDeviation, a.Stupid.Mean, a.Stupid.StandardDeviation, a.Change, rec))
	}
	return r
}

////////////////////////////////////////////////////////////////////////////////
// Segment
////////////////////////////////////////////////////////////////////////////////

func ConvertSegmentAnalysisResults(dir, output string, analysis []Analysis) {
	var str []string

	index := 5
	middle := 4
	ring := 3
	pinky := 2
	palm := 1
	thumb := 0

	indexStrings := getAnalysisStringsSegment("Index Finger", index, analysis)
	middleStrings := getAnalysisStringsSegment("Middle Finger", middle, analysis)
	ringStrings := getAnalysisStringsSegment("Ring Finger", ring, analysis)
	pinkyStrings := getAnalysisStringsSegment("Pinky Finger", pinky, analysis)
	palmStrings := getAnalysisStringsSegment("Palm", palm, analysis)
	thumbStrings := getAnalysisStringsSegment("Thumb", thumb, analysis)

	str = append(str, "############################################################\n")
	str = append(str, "##### Index Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range indexStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Middle Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range middleStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Ring Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range ringStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Pinky Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range pinkyStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Palm Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range palmStrings {
		str = append(str, v)
	}
	str = append(str, "############################################################\n")
	str = append(str, "##### Thumb Finger\n")
	str = append(str, "############################################################\n")
	for _, v := range thumbStrings {
		str = append(str, v)
	}

	WriteStrings(fmt.Sprintf("%s/%s", dir, output), str)
}

func getAnalysisStringsSegment(label string, index int, analysis []Analysis) []string {
	var r []string

	xxIndex := index*9 + 0
	xyIndex := index*9 + 1
	xzIndex := index*9 + 2

	yxIndex := index*9 + 3
	yyIndex := index*9 + 4
	yzIndex := index*9 + 5

	zxIndex := index*9 + 6
	zyIndex := index*9 + 7
	zzIndex := index*9 + 8

	for _, a := range analysis {
		r = appendGoodMCStringSegment(xxIndex, label, "x vs. x", a, r)
		r = appendGoodMCStringSegment(xyIndex, label, "x vs. y", a, r)
		r = appendGoodMCStringSegment(xzIndex, label, "x vs. z", a, r)

		r = appendGoodMCStringSegment(yxIndex, label, "y vs. x", a, r)
		r = appendGoodMCStringSegment(yyIndex, label, "y vs. y", a, r)
		r = appendGoodMCStringSegment(yzIndex, label, "y vs. z", a, r)

		r = appendGoodMCStringSegment(zxIndex, label, "z vs. x", a, r)
		r = appendGoodMCStringSegment(zyIndex, label, "z vs. y", a, r)
		r = appendGoodMCStringSegment(zzIndex, label, "z vs. z", a, r)
	}

	for _, a := range analysis {
		r = appendBadMCStringSegment(xxIndex, label, "x vs. x", a, r)
		r = appendBadMCStringSegment(xyIndex, label, "x vs. y", a, r)
		r = appendBadMCStringSegment(xzIndex, label, "x vs. z", a, r)

		r = appendBadMCStringSegment(yxIndex, label, "y vs. x", a, r)
		r = appendBadMCStringSegment(yyIndex, label, "y vs. y", a, r)
		r = appendBadMCStringSegment(yzIndex, label, "y vs. z", a, r)

		r = appendBadMCStringSegment(zxIndex, label, "z vs. x", a, r)
		r = appendBadMCStringSegment(zyIndex, label, "z vs. y", a, r)
		r = appendBadMCStringSegment(zzIndex, label, "z vs. z", a, r)
	}

	for _, a := range analysis {
		r = appendChangeStringSegment(xxIndex, label, "x vs. x", a, r)
		r = appendChangeStringSegment(xyIndex, label, "x vs. y", a, r)
		r = appendChangeStringSegment(xzIndex, label, "x vs. z", a, r)

		r = appendChangeStringSegment(yxIndex, label, "y vs. x", a, r)
		r = appendChangeStringSegment(yyIndex, label, "y vs. y", a, r)
		r = appendChangeStringSegment(yzIndex, label, "y vs. z", a, r)

		r = appendChangeStringSegment(zxIndex, label, "z vs. x", a, r)
		r = appendChangeStringSegment(zyIndex, label, "z vs. y", a, r)
		r = appendChangeStringSegment(zzIndex, label, "z vs. z", a, r)
	}

	return r
}

func appendGoodMCStringSegment(index int, label string, vs string, a Analysis, r []string) []string {
	if a.Index == index && a.GoodMC == true {
		return append(r, fmt.Sprintf("%s (%s): Good MC with value %.3f\n", label, vs, a.Intelligent.Mean))
	}
	return r
}

func appendBadMCStringSegment(index int, label string, vs string, a Analysis, r []string) []string {
	if a.Index == index && a.BadMC == true {
		return append(r, fmt.Sprintf("%s (%s): Bad MC with value %.3f\n", label, vs, a.Intelligent.Mean))
	}
	return r
}

func appendChangeStringSegment(index int, label string, vs string, a Analysis, r []string) []string {
	if a.Index == index && a.UseChange == true {
		rec := "increase stiffness"
		if a.Change < 0.0 {
			rec = "increase compliance"
		}
		return append(r, fmt.Sprintf("%s (%s): Stable values where detected for Intelligent (mean %.3f, std %.3f) Stupid (mean %.3f, std %.3f). Difference (mean values) is %.3f. Recommendation is %s\n", label, vs, a.Intelligent.Mean, a.Intelligent.StandardDeviation, a.Stupid.Mean, a.Stupid.StandardDeviation, a.Change, rec))
	}
	return r
}
