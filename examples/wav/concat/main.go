package main

import (
	"fmt"
	"os"

	"github.com/liamcr/wavy/cmd/wav"
)

const currentDirectory = "examples/wav/concat"

func main() {
	firstWavFile, err := os.Open(fmt.Sprintf("%s/input0.wav", currentDirectory))
	if err != nil {
		panic(fmt.Sprintf("opening wav file: %v", err.Error()))
	}
	defer func() {
		if err := firstWavFile.Close(); err != nil {
			panic(fmt.Sprintf("closing file: %v", err.Error()))
		}
	}()

	secondWavFile, err := os.Open(fmt.Sprintf("%s/input1.wav", currentDirectory))
	if err != nil {
		panic(fmt.Sprintf("opening wav file: %v", err.Error()))
	}
	defer func() {
		if err := secondWavFile.Close(); err != nil {
			panic(fmt.Sprintf("closing file: %v", err.Error()))
		}
	}()

	firstWav, err := wav.Decode(firstWavFile)
	if err != nil {
		panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
	}

	secondWav, err := wav.Decode(secondWavFile)
	if err != nil {
		panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
	}

	err = firstWav.Concat(secondWav)
	if err != nil {
		panic(fmt.Sprintf("concatenating wav files: %v", err.Error()))
	}

	err = firstWav.Write(fmt.Sprintf("%s/output.wav", currentDirectory))
	if err != nil {
		panic(fmt.Sprintf("Saving concatenated wav file: %v", err.Error()))
	}

	fmt.Printf("Concatenated file saved at %s/output.wav", currentDirectory)
}