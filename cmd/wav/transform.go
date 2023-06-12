package wav

import (
	"errors"
	"fmt"
	"math"
)

// SpeedUp speeds up the wav file by a specified factor.
// This will not only increase the speed of the audio, but
// will increase the pitch as well.
func (w *Wav) SpeedUp(factor float32) {
	w.SampleRate = uint32(factor * float32(w.SampleRate))
}

// SlowDown slows up the wav file by a specified factor.
// This will not only decrease the speed of the audio, but
// will lower the pitch as well.
func (w *Wav) SlowDown(factor float32) {
	w.SampleRate = uint32(float32(w.SampleRate) / factor)
}

// ConvertToMono takes a mono wav struct and converts it to an
// equivalent stereo wav struct
func (w *Wav) ConvertToMono() error {
	if (w.Channels != uint16(2)) {
		return fmt.Errorf("input must have 2 audio channels, but this one has %v", w.Channels)
	}

	w.Channels = 1
	w.DataSize /= 2
	for i := 0; i < len(w.Data); i++ {
		if len(w.Data[i].ChannelData) != 2 {
			return errors.New("malformed wav struct")
		}
		if w.BitsPerSample == uint16(8) {
			firstChannelVal, ok := w.Data[i].ChannelData[0].(uint8)
			if !ok {
				return errors.New("malformed wav struct")
			}

			secondChannelVal, ok := w.Data[i].ChannelData[1].(uint8)
			if !ok {
				return errors.New("malformed wav struct")
			}

			w.Data[i].ChannelData[0] = (firstChannelVal / 2) + (secondChannelVal / 2)
			// Remove lingering second channel data
			w.Data[i].ChannelData = w.Data[i].ChannelData[:1]
		}
		if w.BitsPerSample == uint16(16) {
			firstChannelVal, ok := w.Data[i].ChannelData[0].(int16)
			if !ok {
				return errors.New("malformed wav struct")
			}

			secondChannelVal, ok := w.Data[i].ChannelData[1].(int16)
			if !ok {
				return errors.New("malformed wav struct")
			}

			w.Data[i].ChannelData[0] = (firstChannelVal / 2) + (secondChannelVal / 2)
			// Remove lingering second channel data
			w.Data[i].ChannelData = w.Data[i].ChannelData[:1]
		}
		if w.BitsPerSample == uint16(32) {
			firstChannelVal, ok := w.Data[i].ChannelData[0].(int32)
			if !ok {
				return errors.New("malformed wav struct")
			}

			secondChannelVal, ok := w.Data[i].ChannelData[1].(int32)
			if !ok {
				return errors.New("malformed wav struct")
			}

			w.Data[i].ChannelData[0] = (firstChannelVal / 2) + (secondChannelVal / 2)
			// Remove lingering second channel data
			w.Data[i].ChannelData = w.Data[i].ChannelData[:1]
		}
		if w.BitsPerSample == uint16(64) {
			firstChannelVal, ok := w.Data[i].ChannelData[0].(int64)
			if !ok {
				return errors.New("malformed wav struct")
			}

			secondChannelVal, ok := w.Data[i].ChannelData[1].(int64)
			if !ok {
				return errors.New("malformed wav struct")
			}

			w.Data[i].ChannelData[0] = (firstChannelVal / 2) + (secondChannelVal / 2)
			// Remove lingering second channel data
			w.Data[i].ChannelData = w.Data[i].ChannelData[:1]
		}
	}

	return nil
}

// ConvertToStereo takes a mono wav struct and converts it to an
// equivalent stereo wav struct
func (w *Wav) ConvertToStereo() error {
	if (w.Channels != uint16(1)) {
		return fmt.Errorf("input must have 1 audio channel, but this one has %v", w.Channels)
	}

	w.Channels = 2
	if w.DataSize > math.MaxUint32 / 2 {
		return errors.New("file size too large to be converted to stereo")
	}
	w.DataSize *= 2
	for i := 0; i < len(w.Data); i++ {
		if len(w.Data[i].ChannelData) != 1 {
			return errors.New("malformed wav struct")
		}
		w.Data[i].ChannelData = append(w.Data[i].ChannelData, w.Data[i].ChannelData[0])
	}

	return nil
}