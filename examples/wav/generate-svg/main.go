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

	_, err = newWav.GenerateSvg("output/output.svg", 350, 0, "#ff0")
	if err != nil {
		panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
	}
}