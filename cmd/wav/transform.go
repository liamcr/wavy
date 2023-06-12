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