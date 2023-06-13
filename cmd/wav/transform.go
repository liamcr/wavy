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

	// TODO: Just like we're standardizing the number of channels,
	// we need to think of the case where there are differing sample rates,
	// and differing bit depths
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

	w.Data = append(w.Data, toAdd.Data...)
	w.DataSize = w.DataSize + toAdd.DataSize

	if revertAddedWav {
		err := toAdd.ConvertToMono()
		if err != nil {
			return err
		}
	}

	return nil
}