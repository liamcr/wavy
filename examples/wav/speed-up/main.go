package main

import (
	"fmt"
	"os"

	"github.com/liamcr/wavy/cmd/wav"
)

const currentDirectory = "examples/wav/speed-up"

func main() {
	wavFile, err := os.Open(fmt.Sprintf("%s/input.wav", currentDirectory))
	if err != nil {
		panic(fmt.Sprintf("opening wav file: %v", err.Error()))
	}
	defer func() {
		if err := wavFile.Close(); err != nil {
			panic(fmt.Sprintf("closing file: %v", err.Error()))
		}
	}()

	newWav, err := wav.Decode(wavFile)
	if err != nil {
		panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
	}

	err = newWav.SpeedUp(2)
	if err != nil {
		panic(fmt.Sprintf("speeding up wav file: %v", err.Error()))
	}

	err = newWav.Write(fmt.Sprintf("%s/output.wav", currentDirectory))
	if err != nil {
		panic(fmt.Sprintf("Saving sped-up wav file: %v", err.Error()))
	}

	fmt.Printf("Sped-up file saved at %s/output.wav", currentDirectory)
}