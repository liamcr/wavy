# Wavy

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

### Slow Down

The `.SlowDown` function lets you slow down a wav file by some set factor.
For example, if you set the factor to 2, the resulting wav will be twice
as slow (i.e. the length will be twice as long).

```go
speedyWav.SlowDown(4)

// speedyWav will now be 4 times slower
```
