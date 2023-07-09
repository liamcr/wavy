package main

import (
	"fmt"
	"os"

	"github.com/liamcr/wavy/cmd/wav"
)

func main() {
	wavFile, err := os.Open("input/input.wav")
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

	newWav.SlowDown(2)
	if err != nil {
		panic(fmt.Sprintf("speeding up wav file: %v", err.Error()))
	}

	err = newWav.Write("output/output-slow-down.wav")
	if err != nil {
		panic(fmt.Sprintf("Saving sped-up wav file: %v", err.Error()))
	}
}