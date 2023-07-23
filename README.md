# Wavy

[![Go Reference](https://pkg.go.dev/badge/github.com/liamcr/wavy.svg)](https://pkg.go.dev/github.com/liamcr/wavy)

A Go Library For Transforming Audio Files

## Reading in Wav Files

Reading in a .wav audio file can be done with the `wav.Decode` function:

```go
wavFile, err := os.Open("input/input.wav")
if err != nil {
    panic(fmt.Sprintf("opening wav file: %v", err.Error()))
}

decodedWav, err := wav.Decode(wavFile)
if err != nil {
    panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
}

// Transformations can now be applied to `decodedWav`
```

## Writing Wav Files

`Wav` structs have a handy `.Write()` function to easily write a transformed
wav to a file.

```go
err = transformedWav.Write("output.wav")
if err != nil {
    panic(fmt.Sprintf("Saving wav file: %v", err.Error()))
}
```

## Transformations

There are several transformations that can be applied to wav files.

### Speed Up

The `.SpeedUp` function lets you speed up a wav file by some set factor.
For example, if you set the factor to 2, the resulting wav will be twice
as fast (i.e. the length will be half as long).

```go
err := slowWav.SpeedUp(4)
if err != nil {
    panic(fmt.Sprintf("Speeding up wav file: %v", err.Error()))
}

// slowWav will now be 4 times faster
```

To see an example, run `go run ./examples/wav/speed-up`

### Slow Down

The `.SlowDown` function lets you slow down a wav file by some set factor.
For example, if you set the factor to 2, the resulting wav will be twice
as slow (i.e. the length will be twice as long).

```go
speedyWav.SlowDown(4)

// speedyWav will now be 4 times slower
```

To see an example, run `go run ./examples/wav/slow-down`

### Concatenate Two Audio Files

After importing and decoding two audio files, you can concatenate them together
by using the `.Concat` function.

You can concatenate mono and stereo files together.

```go
err := firstWav.Concat(secondWav)
if err != nil {
    panic(fmt.Sprintf("Concatenating wav files: %v", err.Error()))
}

// firstWav is now the concatenated wav file. secondWav is unaffected.
```

To see an example, run `go run ./examples/wav/concat`

## Convert

### Convert to Mono

Convert a stereo audio file to a mono audio file using the `ConvertToMono` function.

```go
err := stereoWav.ConvertToMono()
if err != nil {
    panic(fmt.Sprintf("Converting wav file: %v", err.Error()))
}

// stereoWav is now a mono audio file
```

### Convert to Stereo

Convert a mono audio file to a stereo audio file using the `ConvertToStereo` function.

```go
err := monoWav.ConvertToStereo()
if err != nil {
    panic(fmt.Sprintf("Converting wav file: %v", err.Error()))
}

// monoWav is now a stereo audio file
```

### Resample

You can adjust the sample rate of an audio file (without adjusting the length) by calling
the `Resample` function.

```go
err := myWav.Resample(44100)
if err != nil {
    panic(fmt.Sprintf("Resampling wav file: %v", err.Error()))
}

// myWav now has a sample rate of 44100 Hz
```

## Other

### Generate SVG

When incorporating this library into any kind of audio manipulation application,
it's important to provide some kind of visual component to the audio files. The
`GenerateSVG` function takes the given wav file and outputs an SVG representing
it's waveform.

```go
output, err = wavForm.GenerateSvg("output.svg", 350, 0, "#ff0")
if err != nil {
    panic(fmt.Sprintf("decoding wav file: %v", err.Error()))
}

// output is an `os.File` representing the resulting svg file
```

To see an example, run `go run ./examples/wav/generate-svg`
