package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/icza/bitio"
)

type Encoder struct {
	WordLength         uint8
	ReconstructionRate int
	Mode               bool
	N                  int
	Reader             bitio.Reader
	Writer             bitio.Writer
	SeenCharacters     map[uint64]*Tree
}

type Decoder struct {
	WordLength         uint8
	ReconstructionRate int
	Mode               bool
	N                  int
	Reader             bitio.Reader
	Writer             bitio.Writer
	SeenCharacters     map[uint64]*Tree
}

func NewEncoder(length uint8, reconstructionRate, n uint, mode bool, inputFilename, outputFilename string) (*Encoder, *bufio.Reader, *bufio.Writer) {
	var inputFile *os.File
	if inputFilename != "-" {
		var err error
		inputFile, err = os.Open(inputFilename)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", inputFilename, err)
			return nil, nil, nil
		}
	}

	var inputReader *bitio.Reader
	var reader *bufio.Reader
	if inputFile != nil {
		reader = bufio.NewReader(inputFile)
		inputReader = bitio.NewReader(reader)
	} else {
		reader = bufio.NewReader(os.Stdin)
		inputReader = bitio.NewReader(reader)
	}

	var outputFile *os.File
	if outputFilename != "" {
		var err error
		outputFile, err = os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", outputFilename, err)
			return nil, nil, nil
		}
	}

	var outputWriter *bitio.Writer
	var writer *bufio.Writer
	if outputFile != nil {
		writer = bufio.NewWriter(outputFile)
		outputWriter = bitio.NewWriter(writer)
	} else {
		writer = bufio.NewWriter(os.Stdout)
		outputWriter = bitio.NewWriter(writer)
	}

	seenCharacters := make(map[uint64]*Tree)
	encoder := &Encoder{
		WordLength:         length,
		ReconstructionRate: int(math.Pow(2, float64(reconstructionRate))),
		Mode:               mode,
		N:                  int(math.Pow(2, float64(n))),
		Reader:             *inputReader,
		Writer:             *outputWriter,
		SeenCharacters:     seenCharacters,
	}

	encoder.Writer.WriteBits(uint64(length), 4)
	encoder.Writer.WriteBits(uint64(reconstructionRate), 4)
	encoder.Writer.WriteBool(mode)
	encoder.Writer.WriteBits(uint64(n), 4)

	return encoder, reader, writer
}

func (e *Encoder) Encode() error {
	order := int(math.Pow(2, float64(e.WordLength))) * 2
	nyt := NewTree(0, 0, order, nil, nil, nil, true)

	var word uint64
	var err error
	for {
		word, err = e.Reader.ReadBits(uint8(e.WordLength))
		if err == io.EOF {
			nyt = e.ProcessWord(word, nyt)
			break
		}
		if err != nil {
			return err
		}
		nyt = e.ProcessWord(word, nyt)
	}

	return nil
}

func (e *Encoder) ProcessWord(word uint64, nyt *Tree) *Tree {
	if _, ok := e.SeenCharacters[word]; !ok {
		var newCharacterPointer *Tree
		WriteNytIndex(nyt, &e.Writer)
		e.Writer.WriteBits(word, e.WordLength)

		newCharacterPointer, nyt = nyt.ProcessNewCharacter(word)
		e.SeenCharacters[word] = newCharacterPointer
		return nyt
	}

	characterPointer := e.SeenCharacters[word]
	treeIndex, length := characterPointer.GetTreeIndex()
	e.Writer.WriteBits(treeIndex, uint8(length))
	characterPointer.Update()

	return nyt
}

func WriteNytIndex(nyt *Tree, writer *bitio.Writer) {
	nytPath, length := nyt.GetTreeIndex()
	fmt.Printf("%b\n", nytPath)
	writer.WriteBits(nytPath, uint8(length))
}

func NewDecoder(inputFilename, outputFilename string) (*Decoder, *bufio.Reader, *bufio.Writer, error) {
	var inputFile *os.File
	if inputFilename != "-" {
		var err error
		inputFile, err = os.Open(inputFilename)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", inputFilename, err)
			return nil, nil, nil, err
		}
	}

	var inputReader *bitio.Reader
	var reader *bufio.Reader
	if inputFile != nil {
		reader = bufio.NewReader(inputFile)
		inputReader = bitio.NewReader(reader)
	} else {
		reader = bufio.NewReader(os.Stdin)
		inputReader = bitio.NewReader(reader)
	}

	var outputFile *os.File
	if outputFilename != "" {
		var err error
		outputFile, err = os.OpenFile(outputFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			fmt.Printf("Failed to open file %s: %s", outputFilename, err)
			return nil, nil, nil, err
		}
	}

	var outputWriter *bitio.Writer
	var writer *bufio.Writer
	if outputFile != nil {
		writer = bufio.NewWriter(outputFile)
		outputWriter = bitio.NewWriter(writer)
	} else {
		writer = bufio.NewWriter(os.Stdout)
		outputWriter = bitio.NewWriter(writer)
	}

	decoder := &Decoder{
		Reader:         *inputReader,
		Writer:         *outputWriter,
		SeenCharacters: map[uint64]*Tree{},
	}
	wordlength, err := decoder.Reader.ReadBits(4)
	if err != nil {
		return nil, nil, nil, err
	}
	reconstructionRate, err := decoder.Reader.ReadBits(4)
	if err != nil {
		return nil, nil, nil, err
	}
	mode, err := decoder.Reader.ReadBool()
	if err != nil {
		return nil, nil, nil, err
	}
	modeCondition, err := decoder.Reader.ReadBits(4)
	if err != nil {
		return nil, nil, nil, err
	}

	decoder.WordLength = uint8(wordlength)
	decoder.ReconstructionRate = int(reconstructionRate)
	decoder.Mode = mode
	decoder.N = int(modeCondition)

	return decoder, reader, writer, nil
}

func (d *Decoder) Decode() error {
	order := int(math.Pow(2, float64(d.WordLength))) * 2
	root := NewTree(0, 0, order, nil, nil, nil, true)

	word, err := d.Reader.ReadBits(d.WordLength)
	if err != nil {
		return err
	}
	root.ProcessNewCharacter(word)
	d.Writer.WriteBits(word, d.WordLength)

	var bit bool
	for err != io.EOF {
		current := root
		for current.Left != nil && current.Right != nil {
			bit, err = d.Reader.ReadBool()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			if bit {
				current = current.Right
			} else {
				current = current.Left
			}
		}
		if current.Nyt {
			word, err := d.Reader.ReadBits(d.WordLength)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			current.ProcessNewCharacter(word)
			d.Writer.WriteBits(word, d.WordLength)
		} else {
			d.Writer.WriteBits(current.Value, d.WordLength)
			current.Update()
		}
	}

	return nil
}
