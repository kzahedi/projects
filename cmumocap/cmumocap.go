package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/kzahedi/goc3d"
	"github.com/kzahedi/goent/continuous"
	"github.com/kzahedi/goent/continuous/state"
	"github.com/kzahedi/utils"

	"gonum.org/v1/plot/plotter"

	pb "gopkg.in/cheggaaa/pb.v1"
)

// labels="C7,T10,LFWT,LKNE,LANK,LSHO,LWRA,RFWT,RKNE,RANK,RSHO,RWRA,LFHD,LBHD,RBHD,RFHD,LHEE,RHEE"
func main() {

	subjectPtr := flag.Int("s", 38, "subject id")
	trialPtr := flag.Int("t", 4, "trial id")
	labelsSetPtr := flag.Int("l", 0, "labels")
	centreCoordinateFramePtr := flag.String("L", "T10", "labels")
	printLabelsPtr := flag.Bool("v", false, "print all the labels")
	useJerk := flag.Bool("j", false, "use jerk instead of curvature")
	flag.Parse()

	var labelSetLabel string
	var selectedLabels []string
	switch *labelsSetPtr {
	case 0: // all labels
		selectedLabels = []string{"T10", "C7", "LFWT", "LKNE", "LANK", "LSHO", "LELB", "LWRA", "RFWT", "RKNE", "RANK", "RSHO", "RELB", "RWRA", "LFHD", "LBHD", "RBHD", "RFHD", "LHEE", "RHEE"}
		labelSetLabel = "full"
	case 1: // legs only
		selectedLabels = []string{"T10", "LFWT", "LKNE", "LANK", "LHEE", "RFWT", "RKNE", "RANK", "RHEE"} // T10 is for the transformation
		labelSetLabel = "legs"
	case 2: // left leg only
		selectedLabels = []string{"T10", "LFWT", "LKNE", "LANK", "LHEE"} // T10 is for the transformation
		labelSetLabel = "left_leg"
	case 3: // left leg only
		selectedLabels = []string{"T10", "RFWT", "RKNE", "RANK", "RHEE"} // T10 is for the transformation
		labelSetLabel = "right_leg"
	}

	sStr := prefix(*subjectPtr)
	tStr := prefix(*trialPtr)
	directory := fmt.Sprintf("%s_%s_%s", sStr, tStr, labelSetLabel)
	id := fmt.Sprintf("%s_%s", sStr, tStr)
	c3dFile := fmt.Sprintf("%s/%s.c3d", directory, id)
	c3dUrl := fmt.Sprintf("http://mocap.cs.cmu.edu/subjects/%s/%s.c3d", sStr, id)

	mpgFile := fmt.Sprintf("%s/%s.mpg", directory, id)
	mpgUrl := fmt.Sprintf("http://mocap.cs.cmu.edu/subjects/%s/%s.mpg", sStr, id)

	aviFile := fmt.Sprintf("%s/%s.avi", directory, id)
	aviUrl := fmt.Sprintf("http://mocap.cs.cmu.edu/subjects/%s/%s.avi", sStr, id)

	headerFile := fmt.Sprintf("%s/%s_meta.c3d", directory, id)
	resultFile := fmt.Sprintf("%s/%s_results.txt", directory, id)
	csvFile := fmt.Sprintf("%s/%s_results.csv", directory, id)
	wFile := fmt.Sprintf("%s/%s_w.csv", directory, id)
	aFile := fmt.Sprintf("%s/%s_a.csv", directory, id)

	utils.CreateDirectory(directory, false)

	downloadFile(c3dFile, c3dUrl)
	downloadFile(mpgFile, mpgUrl)
	downloadFile(aviFile, aviUrl)

	////////////////////////////////////////////////////////////
	// Read C3D File
	////////////////////////////////////////////////////////////
	header, info, data := goc3d.ReadC3D(c3dFile)
	// fmt.Println(header)

	fmt.Println(fmt.Sprintf("Writing header information to %s\n", headerFile))
	file, _ := os.Create(headerFile)
	defer file.Close()
	file.WriteString(fmt.Sprintf("%s\n", header))

	////////////////////////////////////////////////////////////
	// Extract prefix
	////////////////////////////////////////////////////////////

	prefix := ""

	for _, p := range info.Parameters {
		if p.Name == "LABEL_PREFIXES" {
			prefix = strings.Trim(p.StringData[0], " ")
		}
	}

	fmt.Println(fmt.Sprintf("Removing prefix \"%s\"", prefix))

	var labels []string
	for _, p := range info.Parameters {
		if p.Name == "LABELS" {
			for _, s := range p.StringData {
				s = strings.Replace(s, prefix, "", -1)
				s = strings.Trim(s, " ")
				s = strings.Trim(s, "\t")
				labels = append(labels, s)
			}
		}
	}

	for _, p := range info.Parameters {
		if p.Name == "LABELS2" {
			for _, s := range p.StringData {
				s = strings.Replace(s, prefix, "", -1)
				s = strings.Trim(s, " ")
				s = strings.Trim(s, "\t")
				labels = append(labels, s)
			}
		}
	}

	if *printLabelsPtr == true {
		for _, l := range labels {
			fmt.Println(fmt.Sprintf("\"%s\"", l))
		}
		os.Exit(0)
	}

	var foundLabels []string

	for _, s := range selectedLabels {
		found := false
		for _, t := range labels {
			if s == t {
				foundLabels = append(foundLabels, s)
				found = true
			}
		}
		if found == false {
			fmt.Println(fmt.Sprintf("The label \"%s\" was not found", s))
		}
	}

	if len(foundLabels) == 0 {
		fmt.Printf("The specified labels were not found in the data")
		os.Exit(-1)
	}

	// from here on, we will work with foundLabels, which is the
	// intersection of selectedLabels and labels

	var df dataframe.DataFrame

	fmt.Println("Extracting data")
	bar := pb.StartNew(len(foundLabels))
	var indices []int

	// get the indices from the data
	for _, s := range foundLabels {
		for i, t := range labels {
			if s == t {
				indices = append(indices, i)
			}
		}
	}

	for n, i := range indices {
		l := labels[i]
		ls, ds := getData(i, l, data, *useJerk)
		for j := range ls {
			d := dataframe.New(
				series.New(ds[j], series.Float, ls[j]),
			)
			if (n == 0) && (j == 0) {
				df = d
			} else {
				df = df.CBind(d)
			}
		}
		bar.Increment()
	}
	bar.Finish()

	////////////////////////////////////////////////////////////
	// Export W, A File
	////////////////////////////////////////////////////////////

	exportWFile(wFile, df, foundLabels)
	exportAFile(aFile, df, foundLabels)

	////////////////////////////////////////////////////////////
	// Calculating MC_W
	////////////////////////////////////////////////////////////

	df = coordinateTransformation(df, *centreCoordinateFramePtr)
	df = normaliseDataFrame(df)

	foundLabelsWithoutCentre := remove(*centreCoordinateFramePtr, foundLabels)

	w := extractW(df, foundLabelsWithoutCentre)
	a := extractA(df, foundLabelsWithoutCentre)

	nrOfLabels := len(foundLabelsWithoutCentre)
	n, _ := df.Dims()

	w2w1a1 := make([][]float64, n-1, n-1)
	for i := 0; i < n-1; i++ {
		w2w1a1[i] = make([]float64, nrOfLabels*7, nrOfLabels*7) // x', y', z', x, y, z, a
	}

	// w2,w1
	for row := 0; row < n-1; row++ {
		for col := 0; col < nrOfLabels*3; col++ {
			w2w1a1[row][col] = w[row+1][col]            // w2: x', y', z'
			w2w1a1[row][nrOfLabels*3+col] = w[row][col] // w1: x, y, z
		}
	}

	// a1
	for row := 0; row < n-1; row++ {
		for col := 0; col < nrOfLabels; col++ {
			w2w1a1[row][nrOfLabels*6+col] = a[row][col] // a1: a
		}
	}

	w2Indices := make([]int, nrOfLabels*3, nrOfLabels*3) // x, y, z
	w1Indices := make([]int, nrOfLabels*3, nrOfLabels*3) // x, y, z
	a1Indices := make([]int, nrOfLabels, nrOfLabels)     // a

	index := 0
	for i := 0; i < nrOfLabels*3; i++ {
		w2Indices[i] = index
		index++
	}
	for i := 0; i < nrOfLabels*3; i++ {
		w1Indices[i] = index
		index++
	}
	for i := 0; i < nrOfLabels; i++ {
		a1Indices[i] = index
		index++
	}

	mcw := state.MorphologicalComputationW(w2w1a1, w2Indices, w1Indices, a1Indices, 40, true)
	mcwc := continuous.MorphologicalComputationW(w2w1a1, w2Indices, w1Indices, a1Indices, 40, true)

	// fmt.Println(mcw)
	fmt.Println("Result written to", resultFile)
	file, _ = os.Create(resultFile)
	defer file.Close()
	file.WriteString(fmt.Sprintf("MI_w: %f\n", mcwc))
	file.WriteString(fmt.Sprintf("MI_w: %f\n", mcw))
	file.WriteString(fmt.Sprintf("Number of data points: %d\n", len(data.Points[0])))

	// fullMCW := make([]float64, header.LastFrame, header.LastFrame)
	// for i, v := range mcw {
	// fullMCW[header.FirstFrame+i] = v
	// }

	fmt.Println("Point-wise data written to", csvFile)
	// utils.WriteCsvFloatArray(csvFile, matchDataLength(mcw, header), nil)
	utils.WriteCsvFloatArray(csvFile, mcw, nil)
}

func makePoints(data []float64) plotter.XYs {
	n := len(data)
	pts := make(plotter.XYs, n)
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = data[i]
	}
	return pts
}

func matchDataLength(d []float64, header goc3d.C3DHeader) []float64 {
	var s []float64
	for i := 0; i < header.FirstFrame-1; i++ {
		s = append(s, 0.0)
	}
	for _, v := range d {
		s = append(s, v)
	}
	for len(s) < header.LastFrame {
		s = append(s, 0.0)
	}
	return s
}

func getData(index int, label string, data goc3d.C3DData, useJerk bool) ([]string, [][]float64) {
	if index >= len(data.Points) {
		return []string{}, [][]float64{}
	}

	var labels []string

	labels = append(labels, strings.Trim(label, " ")+".X")
	labels = append(labels, strings.Trim(label, " ")+".Y")
	labels = append(labels, strings.Trim(label, " ")+".Z")
	labels = append(labels, strings.Trim(label, " ")+".A")

	//fmt.Println("Reading trajectory", index)

	points := data.Points[index]

	rdata := make([][]float64, 4, 4)
	for i := 0; i < 4; i++ {
		rdata[i] = make([]float64, len(points), len(points))
	}

	for i, p := range points {
		rdata[0][i] = float64(p.X)
		rdata[1][i] = float64(p.Y)
		rdata[2][i] = float64(p.Z)
	}

	// velocity
	for i := 1; i < len(points); i++ {
		xdist := (rdata[0][i] - rdata[0][i-1])
		ydist := (rdata[1][i] - rdata[1][i-1])
		zdist := (rdata[2][i] - rdata[2][i-1])
		dist := math.Sqrt(xdist*xdist + ydist*ydist + zdist*zdist)

		rdata[3][i] = dist
	}

	// acceleration
	for i := 1; i < len(points); i++ {
		rdata[3][i] = rdata[3][i] - rdata[3][i-1]
	}

	if useJerk == true {
		// third derivative
		for i := 1; i < len(points); i++ {
			rdata[3][i] = rdata[3][i] - rdata[3][i-1]
		}

		// jerk derivative
		for i := 1; i < len(points); i++ {
			rdata[3][i] = rdata[3][i] - rdata[3][i-1]
		}
	}

	return labels, rdata
}

func MinMax(data []float64) (float64, float64) {
	min := data[0]
	max := data[0]

	for _, v := range data {
		if v > max {
			max = v
		}
		if v < min {
			v = min
		}
	}
	return min, max
}

func coordinateTransformation(df dataframe.DataFrame, centre string) dataframe.DataFrame {
	fmt.Println(fmt.Sprintf("Transforming data into the local coordinate system of %s", centre))
	names := df.Names()

	xName := fmt.Sprintf("%s.X", centre)
	yName := fmt.Sprintf("%s.Y", centre)
	zName := fmt.Sprintf("%s.Z", centre)

	xData := df.Col(xName).Float()
	yData := df.Col(yName).Float()
	zData := df.Col(zName).Float()

	var r dataframe.DataFrame
	bar := pb.StartNew(len(names))
	for i, name := range names {
		if strings.Contains(name, centre) {
			continue
		}
		s := df.Col(name)
		f := s.Float()
		var cmp []float64
		if strings.Contains(name, ".X") {
			cmp = xData
		}
		if strings.Contains(name, ".Y") {
			cmp = yData
		}
		if strings.Contains(name, ".Z") {
			cmp = zData
		}
		var d dataframe.DataFrame
		if len(cmp) > 0 {
			for i, _ := range f {
				f[i] = f[i] - cmp[i]
			}
			d = dataframe.New(series.New(f, series.Float, name))
		} else {
			d = dataframe.New(series.New(f, series.Float, name))
		}
		if i == 0 {
			r = d
		} else {
			r = r.CBind(d)
		}
		bar.Increment()
	}
	bar.Finish()
	return r
}

func normaliseDataFrame(df dataframe.DataFrame) dataframe.DataFrame {
	fmt.Println("Normalising data")
	names := df.Names()

	var r dataframe.DataFrame

	bar := pb.StartNew(len(names))
	for i, name := range names {
		s := df.Col(name)
		f := s.Float()
		min, max := MinMax(f)
		if math.Abs(max-min) > 0.0001 {
			for i := range f {
				f[i] = (f[i] - min) / (max - min)
			}
		}
		d := dataframe.New(
			series.New(f, series.Float, name),
		)
		if i == 0 {
			r = d
		} else {
			r = r.CBind(d)
		}
		bar.Increment()
	}
	bar.Finish()
	return r
}

func downloadFile(filename, url string) (err error) {

	if utils.FileExists(filename) == true {
		return nil
	}

	fmt.Println(fmt.Sprintf("Downloading %s from %s", filename, url))

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func prefix(n int) string {
	if n < 10 {
		return fmt.Sprintf("0%d", n)
	}
	return fmt.Sprintf("%d", n)
}

func exportWFile(filename string, df dataframe.DataFrame, labels []string) {
	w := extractW(df, labels)

	var header []string

	for _, l := range labels {
		header = append(header, fmt.Sprintf("%s.X", l))
		header = append(header, fmt.Sprintf("%s.Y", l))
		header = append(header, fmt.Sprintf("%s.Z", l))
	}

	fmt.Println("Writing world states to", filename)
	utils.WriteCsvFloatMatrix(filename, w, header)
}

func exportAFile(filename string, df dataframe.DataFrame, labels []string) {
	a := extractA(df, labels)

	var header []string

	for _, l := range labels {
		header = append(header, fmt.Sprintf("%s.X", l))
		header = append(header, fmt.Sprintf("%s.Y", l))
		header = append(header, fmt.Sprintf("%s.Z", l))
	}

	fmt.Println("Writing actuator states to", filename)
	utils.WriteCsvFloatMatrix(filename, a, header)
}

func extractW(df dataframe.DataFrame, labels []string) [][]float64 {
	nrOfLabels := len(labels)

	n, _ := df.Dims()
	w := make([][]float64, n, n)

	for r := 0; r < n; r++ {
		w[r] = make([]float64, 3*nrOfLabels, 3*nrOfLabels)
	}

	var header []string

	for i, label := range labels {
		xName := fmt.Sprintf("%s.X", label)
		yName := fmt.Sprintf("%s.Y", label)
		zName := fmt.Sprintf("%s.Z", label)

		header = append(header, xName)
		header = append(header, yName)
		header = append(header, zName)

		xData := df.Col(xName).Float()
		yData := df.Col(yName).Float()
		zData := df.Col(zName).Float()

		for r := 0; r < n; r++ {
			w[r][i*3+0] = xData[r]
			w[r][i*3+1] = yData[r]
			w[r][i*3+2] = zData[r]
		}
	}

	return w
}

func extractA(df dataframe.DataFrame, labels []string) [][]float64 {
	nrOfLabels := len(labels)

	n, _ := df.Dims()
	a := make([][]float64, n, n)

	for r := 0; r < n; r++ {
		a[r] = make([]float64, nrOfLabels, nrOfLabels)
	}

	var header []string

	for i, label := range labels {
		name := fmt.Sprintf("%s.A", label)
		header = append(header, name)

		data := df.Col(name).Float()

		for r := 0; r < n; r++ {
			a[r][i] = data[r]
		}
	}

	return a
}

func remove(elem string, list []string) []string {
	var r []string
	for _, l := range list {
		if l != elem {
			r = append(r, l)
		}
	}
	return r
}
