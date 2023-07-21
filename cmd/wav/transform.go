package wav

import (
	"errors"
	"math"
)

// SpeedUp speeds up the wav file by a specified factor.
// This will not only increase the speed of the audio, but
// will increase the pitch as well.
func (w *Wav) SpeedUp(factor float64) error {
	if w.SampleRate > math.MaxUint32 / uint32(math.Ceil(factor)) {
		return errors.New("resulting sample rate would be too large (> max uint32)")
	}
	w.SampleRate = uint32(factor * float64(w.SampleRate))
	return nil
}

// SlowDown slows up the wav file by a specified factor.
// This will not only decrease the speed of the audio, but
// will lower the pitch as well.
func (w *Wav) SlowDown(factor float32) {
	w.SampleRate = uint32(float32(w.SampleRate) / factor)
}

// Concat takes another Wav struct and stitches the two audio files
// together. The returned wav struct will have the number of channels
// equal to the largest number of channels out of the two wavs being
// concatenated.
func (w *Wav) Concat(toAdd *Wav) error {
	revertAddedWav := false

	if w.Channels == 2 || toAdd.Channels == 2 {
		if w.Channels == 1 {
			err := w.ConvertToStereo()
			if err != nil {
				return err
			}
		} else if toAdd.Channels == 1 {
			err := toAdd.ConvertToStereo()
			if err != nil {
				return err
			}

			// So as to not cause any side effects to the added wav,
			// set this flag to remember that the `toAdd` wav should
			// be converted back to mono.
			// Important to note that going from mono to stereo, then back
			// to mono results in no data loss.
			revertAddedWav = true
		}
	}

	if w.SampleRate != toAdd.SampleRate {
		// Resample the current audio file so as to not introduce side effects
		// to the file being appended.
		w.Resample(toAdd.SampleRate)
	}

	w.Data = append(w.Data, toAdd.Data...)
	maxBitDepth := uint16(math.Max(float64(w.BitsPerSample), float64(toAdd.BitsPerSample)))

	// If there are differing bit depths we should normalize them
	if w.BitsPerSample != toAdd.BitsPerSample {
		for _, sampleGroup := range(w.Data) {
			for j := 0; j < int(w.Channels); j++ {
				intVal, err := CastToInt(sampleGroup.ChannelData[j])
				if err != nil {
					return err
				}
				if maxBitDepth == uint16(16) {
					sampleGroup.ChannelData[j] = int16(intVal)
				}
				if maxBitDepth == uint16(32) {
					sampleGroup.ChannelData[j] = int32(intVal)
				}
				if maxBitDepth == uint16(64) {
					sampleGroup.ChannelData[j] = int64(intVal)
				}
			}
		}
	}

	w.DataSize = uint32(len(w.Data) + len(toAdd.Data) * int(w.Channels) * (int(maxBitDepth) / 8))

	if revertAddedWav {
		err := toAdd.ConvertToMono()
		if err != nil {
			return err
		}
	}

	return nil
}