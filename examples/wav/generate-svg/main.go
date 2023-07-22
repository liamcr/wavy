package main

import (
	"fmt"
	"os"

	"github.com/liamcr/wavy/cmd/wav"
)

const currentDirectory = "examples/wav/generate-svg"

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

	_, err = newWav.GenerateSvg(fmt.Sprintf("%s/output.svg", currentDirectory), 600, 0, "#4287f5")
	if err != nil {
		panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
	}

	fmt.Printf("Waveform SVG saved at %s/output.svg", currentDirectory)
}