package wav

import (
	"errors"
	"fmt"
	"math"
)

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

// Resample will update the sample rate of the wave file, without
// changing the duration or pitch.
func (w *Wav) Resample(newSampleRate uint32) error {
	sampleDifferenceRatio := float64(w.SampleRate) / float64(newSampleRate)
	newData := []SampleGroup{}

	rawChannelData := make([][]int, int(w.Channels))

	for i := 0; i < len(w.Data); i++ {
		for j := 0; j < int(w.Channels); j++ {
			intValue, err := CastToInt(w.Data[i].ChannelData[j])
			if err != nil {
				return err
			}
			rawChannelData[j] = append(rawChannelData[j], intValue)
		}
	}

	for i := 0.0; i < float64(len(w.Data)); i += sampleDifferenceRatio {
		sampleGroupToAdd := SampleGroup{ChannelData: []any{}}

		for j := 0; j < int(w.Channels); j++ {
			sampleGroupToAdd.ChannelData = append(
				sampleGroupToAdd.ChannelData,
				resamplePoint(
					i,
					rawChannelData[j],
					w.SampleRate,
					1000,
					4.0,
				),
			)
		}

		newData = append(newData, sampleGroupToAdd)
	}

	// covert data back to format that it was previously in (int16, for example)
	finalDataArray := []SampleGroup{}
	for i := 0; i < len(newData); i++ {
		newSampleGroup := SampleGroup{
			ChannelData: []any{},
		}
		for j := 0; j < int(w.Channels); j++ {
			intVal, ok := newData[i].ChannelData[j].(int)
			if !ok {
				return fmt.Errorf("could not convert %v to int", newData[i].ChannelData[j])
			}

			if w.BitsPerSample == uint16(8) {
				newSampleGroup.ChannelData = append(newSampleGroup.ChannelData, uint8(intVal))
			}
			if w.BitsPerSample == uint16(16) {
				newSampleGroup.ChannelData = append(newSampleGroup.ChannelData, int16(intVal))
			}
			if w.BitsPerSample == uint16(32) {
				newSampleGroup.ChannelData = append(newSampleGroup.ChannelData, int32(intVal))
			}
			if w.BitsPerSample == uint16(64) {
				newSampleGroup.ChannelData = append(newSampleGroup.ChannelData, int64(intVal))
			}
		}
		finalDataArray = append(finalDataArray, newSampleGroup)
	}
	
	w.Data = finalDataArray
	w.SampleRate = uint32(newSampleRate)

	bytesPerSample := uint32(w.BitsPerSample) / 8
	w.DataSize = bytesPerSample * uint32(len(w.Data)) * uint32(w.Channels)

	return nil
}

func resamplePoint(
	x float64,
	initialSamples []int,
	initialSampleRate uint32,
	filterCutoffFrequency float64,
	windowWidth float64,
) int {
	var r_w, r_a, sincVal float64

	gainCorrectionFactor := 2 * filterCutoffFrequency / float64(initialSampleRate)
	filteredSample := 0.0

	for i := -windowWidth / 2; i < windowWidth/2; i++ {
		inputSampleIndex := int(x + float64(i))

		// TODO: Rename these variables
		r_w = 0.5 - 0.5 * math.Cos(2 * math.Pi * (0.5 + (float64(inputSampleIndex) - x) / windowWidth))
		r_a = 2 * math.Pi * (float64(inputSampleIndex) - x) * filterCutoffFrequency / float64(initialSampleRate)
		sincVal = 1

		if r_a != 0 {
			sincVal = math.Sin(r_a) / r_a
		}

		if inputSampleIndex >= 0 && inputSampleIndex < len(initialSamples) {
			filteredSample += gainCorrectionFactor * r_w * sincVal * float64(initialSamples[inputSampleIndex])
		}
	}

	return int(filteredSample)
}