package main

import (
	"bytes"
	"testing"

	"github.com/icza/bitio"
)

func TestEncodeDecodeWithDummyInput(t *testing.T) {
	// Dummy input
	dummyInput := "Hello World!"
	inputBuffer := bytes.NewBufferString(dummyInput)
	encodedBuffer := &bytes.Buffer{}

	seenChars := make(map[uint64]*Tree)
	encoder := &Encoder{
		WordLength:     8, // set these values as per your requirements
		Reader:         *bitio.NewReader(inputBuffer),
		Writer:         *bitio.NewWriter(encodedBuffer),
		SeenCharacters: seenChars,
	}
	defer encoder.Writer.Close()
	if err := encoder.Encode(); err != nil {
		t.Errorf("Failed to encode: %s", err)
	}

	decodedBuffer := bytes.NewBuffer(encodedBuffer.Bytes())
	decodedOutputBuffer := new(bytes.Buffer)

	seenCharsDecoder := make(map[uint64]*Tree)
	decoder := &Decoder{
		WordLength:     8, // set these values as per your requirements
		Reader:         *bitio.NewReader(decodedBuffer),
		Writer:         *bitio.NewWriter(decodedOutputBuffer), // Assuming you don't need the decoded output for this test
		SeenCharacters: seenCharsDecoder,
	}
	defer decoder.Writer.Close()

	// Perform decoding
	if err := decoder.Decode(); err != nil {
		t.Errorf("Failed to decode: %s", err)
	}
	decodedOutput := decodedOutputBuffer.String()
	if decodedOutput != dummyInput {
		t.Errorf("Decoded output does not match original input. Got: %s, Want: %s", decodedOutput, dummyInput)
	}

	// Here you can verify the decoded output
	// Note: This step depends on how you can access the decoded data.
	// Compare it with `dummyInput` if necessary.
}
