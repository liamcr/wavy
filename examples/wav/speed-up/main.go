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

	newWav, err := wav.Decode(wavFile)
	if err != nil {
		panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
	}

	newWav.SpeedUp(2)

	err = newWav.Write("output/output-speed-up.wav")
	if err != nil {
		panic(fmt.Sprintf("Saving sped-up wav file: %v", err.Error()))
	}
}