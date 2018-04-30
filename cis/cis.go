package main

import (
	"flag"
	"fmt"
	"os"

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

	a, _ := utils.ReadFloatCsv(aFile)
	w, _ := utils.ReadFloatCsv(wFile)
	l, _ := utils.ReadCsv(lFile)

	fmt.Println(a)
	fmt.Println(w)
	fmt.Println(l)

}
