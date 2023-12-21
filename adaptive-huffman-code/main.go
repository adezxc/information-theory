package main

import (
	"flag"
)

func main() {
	var d bool
	var m uint
	var mode bool
	var k uint
	var n uint
	var inputFlag string
	var outputFlag string
	flag.BoolVar(&d, "decode", false, "If set to true, tries to decode the file.")
	flag.UintVar(&m, "length", 8, "Word length in bits. From 2 to 16")
	flag.BoolVar(&mode, "mode", false, "Mode of the algorithm. If set to true, "+
		"zeroes all the word frequencies after reaching 2^n. If set to false, stops"+
		" reconstruction after reaching 2^n. Defaults to false.")
	flag.UintVar(&k, "freq", 0, "Power of two, on how often to reconstruct the tree. From 0 to 15")
	flag.UintVar(&n, "n", 0, "Power of two, when to trigger mode condition, see --mode. From 0 to 15")
	flag.StringVar(&inputFlag, "input", "-", "If not defined, read from STDIN")
	flag.StringVar(&outputFlag, "output", "", "If not defined, output to STDOUT")
	flag.Parse()

	if d {
		decoder, _, writer, err := NewDecoder(inputFlag, outputFlag)
		if err != nil {
			defer writer.Flush()
			panic(err)
		}
		decoder.Decode()
		defer writer.Flush()
	} else {
		encoder, _, writer := NewEncoder(uint8(m), k, n, mode, inputFlag, outputFlag)
		encoder.Encode()
		defer writer.Flush()
	}

	return
}
