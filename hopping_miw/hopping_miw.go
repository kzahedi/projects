package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kzahedi/goent/continuous"
	cstate "github.com/kzahedi/goent/continuous/state"
	"github.com/kzahedi/goent/dh"
	"github.com/kzahedi/goent/discrete"
	dstate "github.com/kzahedi/goent/discrete/state"
)

func main() {

	files := [3]string{"musfib.csv", "muslin.csv", "dcmot.csv"}
	url := "https://raw.githubusercontent.com/kzahedi/entropy/master/experiments/hopping/data/"
	k := 30
	bins := 300

	for _, file := range files {
		downloadFile(file, url)
	}

	////////////////////////////////////////////////////////////
	// loading all data
	////////////////////////////////////////////////////////////

	musfibData := readDataRaw("musfib.csv")
	musfibPosition := getColumn(musfibData, 1)
	musfibVelocity := getColumn(musfibData, 2)
	musfibAcceleration := getColumn(musfibData, 3)
	// musfibSensor := getColumn(musfibData, 4)
	musfibAction := getColumn(musfibData, 9)

	muslinData := readDataRaw("muslin.csv")
	muslinPosition := getColumn(muslinData, 1)
	muslinVelocity := getColumn(muslinData, 2)
	muslinAcceleration := getColumn(muslinData, 3)
	// muslinSensor := getColumn(muslinData, 4)
	muslinAction := getColumn(muslinData, 9)

	dcmotData := readDataRaw("dcmot.csv")
	dcmotPosition := getColumn(dcmotData, 1)
	dcmotVelocity := getColumn(dcmotData, 2)
	dcmotAcceleration := getColumn(dcmotData, 3)
	// dcmotSensor := combine(dcmotPosition, dcmotVelocity, nil)
	dcmotAction := getColumn(dcmotData, 9)

	////////////////////////////////////////////////////////////
	// discretising for discrete measures
	////////////////////////////////////////////////////////////

	// min / max values
	pos_min, pos_max := getMinMaxValues(musfibPosition, muslinPosition, dcmotPosition)
	vel_min, vel_max := getMinMaxValues(musfibVelocity, muslinVelocity, dcmotVelocity)
	acc_min, acc_max := getMinMaxValues(musfibAcceleration, muslinAcceleration, dcmotAcceleration)
	// musfib and muslin are already normalised, therefore we only get min max for dcmot
	act_min, act_max := getMinMaxValues(dcmotAction, dcmotAction, dcmotAction)
	// dcmot sensors is position & velocity
	// sen_min, sen_max := getMinMaxValues(musfibSensor, muslinSensor, muslinSensor)
	// dc_pos_min, dc_pos_max := getMinMaxValues(dcmotPosition, dcmotPosition, dcmotPosition)
	// dc_vel_min, dc_vel_max := getMinMaxValues(dcmotVelocity, dcmotVelocity, dcmotVelocity)

	// discretising per data stream
	musfibDiscretePosition := dh.DiscretiseVector(musfibPosition, bins, pos_min, pos_max)
	musfibDiscreteVelocity := dh.DiscretiseVector(musfibVelocity, bins, vel_min, vel_max)
	musfibDiscreteAcceleration := dh.DiscretiseVector(musfibAcceleration, bins, acc_min, acc_max)
	musfibDiscreteAction := dh.DiscretiseVector(musfibAction, bins, -1.0, 1.0)
	// musfibDiscreteSensor := dh.DiscretiseVector(musfibSensor, bins, sen_min, sen_max)

	muslinDiscretePosition := dh.DiscretiseVector(muslinPosition, bins, pos_min, pos_max)
	muslinDiscreteVelocity := dh.DiscretiseVector(muslinVelocity, bins, vel_min, vel_max)
	muslinDiscreteAcceleration := dh.DiscretiseVector(muslinAcceleration, bins, acc_min, acc_max)
	muslinDiscreteAction := dh.DiscretiseVector(muslinAction, bins, -1.0, 1.0)
	// muslinDiscreteSensor := dh.DiscretiseVector(muslinSensor, bins, sen_min, sen_max)

	dcmotDiscretePosition := dh.DiscretiseVector(dcmotPosition, bins, pos_min, pos_max)
	dcmotDiscreteVelocity := dh.DiscretiseVector(dcmotVelocity, bins, vel_min, vel_max)
	dcmotDiscreteAcceleration := dh.DiscretiseVector(dcmotAcceleration, bins, acc_min, acc_max)
	dcmotDiscreteAction := dh.DiscretiseVector(dcmotAction, bins, act_min, act_max)
	// dcmotDiscreteSensor := dh.Discretise(dcmotSensor, []int{bins, bins}, []float64{dc_pos_min, dc_vel_min}, []float64{dc_pos_max, dc_vel_max})

	// creating W data stream

	musfibDiscreteW := combine(musfibDiscretePosition, musfibDiscreteVelocity, musfibDiscreteAcceleration)
	muslinDiscreteW := combine(muslinDiscretePosition, muslinDiscreteVelocity, muslinDiscreteAcceleration)
	dcmotDiscreteW := combine(dcmotDiscretePosition, dcmotDiscreteVelocity, dcmotDiscreteAcceleration)

	musfibW := dh.MakeUnivariateRelabelled(musfibDiscreteW, []int{300, 300, 300})
	muslinW := dh.MakeUnivariateRelabelled(muslinDiscreteW, []int{300, 300, 300})
	dcmotW := dh.MakeUnivariateRelabelled(dcmotDiscreteW, []int{300, 300, 300})

	// create A container
	musfibA := dh.Relabel(musfibDiscreteAction)
	muslinA := dh.Relabel(muslinDiscreteAction)
	dcmotA := dh.Relabel(dcmotDiscreteAction)

	// generate w2,w1,a1 data container
	w2w1a1MusFib := make([][]int, len(musfibDiscretePosition)-1, len(musfibDiscretePosition)-1)
	w2w1a1MusLin := make([][]int, len(musfibDiscretePosition)-1, len(musfibDiscretePosition)-1)
	w2w1a1DcMot := make([][]int, len(musfibDiscretePosition)-1, len(musfibDiscretePosition)-1)

	for row := 0; row < len(musfibDiscretePosition)-1; row++ {
		w2w1a1MusFib[row] = make([]int, 3, 3)
		w2w1a1MusFib[row][0] = musfibW[row+1]
		w2w1a1MusFib[row][1] = musfibW[row]
		w2w1a1MusFib[row][2] = musfibA[row]
	}

	for row := 0; row < len(muslinDiscretePosition)-1; row++ {
		w2w1a1MusLin[row] = make([]int, 3, 3)
		w2w1a1MusLin[row][0] = muslinW[row+1]
		w2w1a1MusLin[row][1] = muslinW[row]
		w2w1a1MusLin[row][2] = muslinA[row]
	}

	for row := 0; row < len(dcmotDiscretePosition)-1; row++ {
		w2w1a1DcMot[row] = make([]int, 3, 3)
		w2w1a1DcMot[row][0] = dcmotW[row+1]
		w2w1a1DcMot[row][1] = dcmotW[row]
		w2w1a1DcMot[row][2] = dcmotA[row]
	}

	pw2w1a1MusFib := discrete.Emperical3D(w2w1a1MusFib)
	pw2w1a1MusLin := discrete.Emperical3D(w2w1a1MusLin)
	pw2w1a1DcMot := discrete.Emperical3D(w2w1a1DcMot)

	fmt.Println(fmt.Sprintf("MusFib MI_W (discrete): %f", musFibMIW))
	fmt.Println(fmt.Sprintf("MusLin MI_W (discrete): %f", musLinMIW))
	fmt.Println(fmt.Sprintf("DC-Mot MI_W (discrete): %f", dcmotMIW))

	musFibMIWsd := dstate.MorphologicalComputationW(w2w1a1MusFib)
	musLinMIWsd := dstate.MorphologicalComputationW(w2w1a1MusLin)
	dcmotMIWsd := dstate.MorphologicalComputationW(w2w1a1DcMot)

	f, err := os.Create("mi_w_averaged_results_discrete.csv")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(fmt.Sprintf("MusFib MI_W (discrete): %f", musFibMIW))
	f.WriteString(fmt.Sprintf("MusLin MI_W (discrete): %f", musLinMIW))
	f.WriteString(fmt.Sprintf("DC-Mot MI_W (discrete): %f", dcmotMIW))

	writeToCSV("musfib_mi_w_sd_discrete.csv", musFibMIWsd)
	writeToCSV("muslin_mi_w_sd_discrete.csv", musLinMIWsd)
	writeToCSV("dcmot_mi_w_sd_discrete.csv", dcmotMIWsd)

	////////////////////////////////////////////////////////////
	// continuous
	////////////////////////////////////////////////////////////
	musFibC := readDataForContinuousMIW("musfib.csv")
	musLinC := readDataForContinuousMIW("muslin.csv")
	dcmotC := readDataForContinuousMIW("dcmot.csv")

	musFibC = continuous.Normalise(musFibC)
	musLinC = continuous.Normalise(musLinC)
	dcmotC = continuous.Normalise(dcmotC)

	w2Indices := []int{0, 1, 2}
	w1Indices := []int{3, 4, 5}
	a1Indices := []int{6}

	musFibCMIW := continuous.MorphologicalComputationW(musFibC, w2Indices, w1Indices, a1Indices, k, true)
	musLinCMIW := continuous.MorphologicalComputationW(musLinC, w2Indices, w1Indices, a1Indices, k, true)
	dcmotCMIW := continuous.MorphologicalComputationW(dcmotC, w2Indices, w1Indices, a1Indices, k, true)

	musFibCMIWsd := cstate.MorphologicalComputationW(musFibC, w2Indices, w1Indices, a1Indices, k, true)
	musLinCMIWsd := cstate.MorphologicalComputationW(musLinC, w2Indices, w1Indices, a1Indices, k, true)
	dcmotCMIWsd := cstate.MorphologicalComputationW(dcmotC, w2Indices, w1Indices, a1Indices, k, true)

	f, err = os.Create("mi_w_averaged_results_continuous.csv")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString(fmt.Sprintf("MusFib MI_W: %f", musFibCMIW))
	f.WriteString(fmt.Sprintf("MusLin MI_W: %f", musLinCMIW))
	f.WriteString(fmt.Sprintf("DC-Mot MI_W: %f", dcmotCMIW))

	fmt.Println(fmt.Sprintf("MusFib MI_W (continuous): %f", musFibCMIW))
	fmt.Println(fmt.Sprintf("MusLin MI_W (continuous): %f", musLinCMIW))
	fmt.Println(fmt.Sprintf("DC-Mot MI_W (continuous): %f", dcmotCMIW))

	writeToCSV("musfib_mi_w_sd_continuous.csv", musFibCMIWsd)
	writeToCSV("muslin_mi_w_sd_continuous.csv", musLinCMIWsd)
	writeToCSV("dcmot_mi_w_sd_continuous.csv", dcmotCMIWsd)
}

func convertToString(d []float64) []string {
	r := make([]string, len(d), len(d))

	for i, v := range d {
		r[i] = fmt.Sprintf("%f", v)
	}
	return r
}

func combine(a, b, c []int) [][]int {
	r := make([][]int, len(a), len(a))
	for row := 0; row < len(a); row++ {
		if c != nil {
			r[row] = make([]int, 3, 3)
			r[row][0] = a[row]
			r[row][1] = b[row]
			r[row][2] = c[row]
		} else {
			r[row] = make([]int, 2, 2)
			r[row][0] = a[row]
			r[row][1] = b[row]
		}
	}
	return r
}

func getMin(a, b, c, d float64) float64 {
	return math.Min(a, math.Min(b, math.Min(c, d)))
}

func getMax(a, b, c, d float64) float64 {
	return math.Max(a, math.Max(b, math.Max(c, d)))
}

func getMinMaxValues(a, b, c []float64) (float64, float64) {
	rows := len(a)
	min := a[0]
	max := a[0]
	for i := 0; i < rows; i++ {
		min = getMin(a[i], b[i], c[i], min)
		max = getMax(a[i], b[i], c[i], max)
	}
	return min, max
}

func getColumn(data [][]float64, col int) []float64 {
	r := make([]float64, len(data), len(data))
	for row := 0; row < len(data); row++ {
		r[row] = data[row][col]
	}
	return r
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func downloadFile(filename, url string) (err error) {

	if exists(filename) == true {
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
	resp, err := http.Get(fmt.Sprintf("%s/%s", url, filename))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func readDataForContinuousMIW(filename string) (r [][]float64) {

	var raw [][]float64

	f, _ := os.Open(filename)
	defer f.Close()
	reader := csv.NewReader(bufio.NewReader(f))

	lineCount := 0

	record, err := reader.Read()
	for err != io.EOF {

		if strings.HasPrefix(record[0], "#") {
			record, err = reader.Read()
			continue
		}

		d := make([]float64, 4, 4)
		w1, _ := strconv.ParseFloat(record[1], 64)
		w2, _ := strconv.ParseFloat(record[2], 64)
		w3, _ := strconv.ParseFloat(record[3], 64)
		a1, _ := strconv.ParseFloat(record[9], 64)

		d[0] = w1
		d[1] = w2
		d[2] = w3
		d[3] = a1

		raw = append(raw, d)

		lineCount++
		// fmt.Print(fmt.Sprintf("Line count: %d\r", lineCount))

		record, err = reader.Read()
	}

	// fmt.Println(fmt.Sprintf("\nRead %d lines from %s", lineCount, filename))

	// fmt.Println("Converting raw data to (w2, w1, a1)")

	for i := 0; i < len(raw)-1; i++ {
		d := make([]float64, 7, 7)
		d[0] = raw[i+1][0]
		d[1] = raw[i+1][1]
		d[2] = raw[i+1][2]
		d[3] = raw[i][0]
		d[4] = raw[i][1]
		d[5] = raw[i][2]
		d[6] = raw[i][3]
		r = append(r, d)
	}

	return
}

func readDataRaw(filename string) [][]float64 {

	f, _ := os.Open(filename)
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	d := make([][]float64, len(records), len(records))
	for i, r := range records {
		if strings.HasPrefix(r[0], "#") {
			continue
		}
		d[i] = make([]float64, len(r), len(r))
		for j := 0; j < len(r); j++ {
			d[i][j], _ = strconv.ParseFloat(r[j], 64)
		}
	}

	return d
}

func writeToCSV(filename string, data []float64) {
	dataStr := convertToString(data)
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	w := csv.NewWriter(f)
	w.Write(dataStr)
}
