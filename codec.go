package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/icza/bitio"
)

type Encoder struct {
	WordLength         int
	ReconstructionRate int
	Mode               bool
	N                  int
	Reader             bitio.Reader
	Writer             bitio.Writer
	FrequencyMap       map[uint64]uint64
}

type Node struct {
	Value     *int
	LeftNode  *Node
	RightNode *Node
}

func NewEncoder(length, reconstructionRate, n int, mode bool, inputFilename, outputFilename string) *Encoder {
	var inputFile *os.File
	if inputFilename != "-" {
		var err error
		inputFile, err = os.Open(inputFilename)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", inputFilename, err)
			return nil
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
	if outputFilename != "" {
		var err error
		outputFile, err = os.OpenFile(outputFilename, os.O_CREATE, os.ModeAppend)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", outputFilename, err)
			return nil
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
	return &Encoder{
		WordLength:         length,
		ReconstructionRate: reconstructionRate,
		Mode:               mode,
		N:                  n,
		Reader:             *inputReader,
		Writer:             *outputWriter,
	}
}

func (e *Encoder) Encode() error {

	for {
		word, err := e.Reader.ReadBits(uint8(e.WordLength))
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		needNewTree := e.UpdateFrequencyMap(word)
		if needNewTree {
			e.RebuildTree()
		}
	}

	return nil
}

func (e *Encoder) Decode() error {
	return nil
}

func (e *Encoder) UpdateFrequencyMap(word uint64) bool {
	var needUpdate bool
	if _, ok := e.FrequencyMap[word]; ok {
		e.FrequencyMap[word]++
		needUpdate = false
	} else {
		e.FrequencyMap[word] = 1
		needUpdate = true
	}

	return needUpdate
}

func (e *Encoder) RebuildTree() {

}
