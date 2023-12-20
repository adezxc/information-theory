package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/icza/bitio"
)

func main() {
	var m int
	var mode bool
	var k int
	var n int
	var inputFlag string
	var outputFlag string
	flag.IntVar(&m, "length", 8, "Word length in bits")
	flag.BoolVar(&mode, "mode", false, "Mode of the algorithm. If set to true, "+
		"zeroes all the word frequencies after reaching 2^n. If set to false, stops"+
		" reconstruction after reaching 2^n. Defaults to false.")
	flag.IntVar(&k, "freq", 0, "Power of two, on how often to reconstruct the tree")
	flag.IntVar(&n, "n", 0, "Power of two, when to trigger mode condition, see --mode.")
	flag.StringVar(&inputFlag, "input", "-", "If not defined, read from STDIN")
	flag.StringVar(&outputFlag, "output", "", "If not defined, output to STDOUT")
	flag.Parse()

	var inputFile *os.File
	if inputFlag != "-" {
		var err error
		inputFile, err = os.Open(inputFlag)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", inputFlag, err)
			return
		}
		fmt.Printf(inputFile.Name())
	}

	var inputReader *bitio.Reader
	if inputFile != nil {
		reader := bufio.NewReader(inputFile)
		inputReader = bitio.NewReader(reader)
	} else {
		reader := bufio.NewReader(os.Stdin)
		inputReader = bitio.NewReader(reader)
	}

	var outputFile *os.File
	if outputFlag != "" {
		var err error
		outputFile, err = os.OpenFile(outputFlag, os.O_CREATE, os.ModeAppend)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", outputFlag, err)
			return
		}
		fmt.Printf(outputFile.Name())
	}

	var outputWriter *bitio.Writer
	if outputFile != nil {
		writer := bufio.NewWriter(outputFile)
		outputWriter = bitio.NewWriter(writer)
	} else {
		writer := bufio.NewWriter(os.Stdout)
		outputWriter = bitio.NewWriter(writer)
	}

	return
}
