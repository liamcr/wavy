package wav

import (
	"fmt"

	"github.com/liamcr/wavy/internal/util"
)

// GenerateBucketedAvgSampleVals separates the audio file data into a set number of
// buckets, by dividing the data array into equal parts, and averaging the amplitudes
// over each part
func (w *Wav) GenerateBucketedAvgSampleVals(buckets, channel int) ([]float64, error) {
	bucketVals := make([]float64, buckets)
	samplesInBuckets := len(w.Data) / buckets

	if channel >= int(w.Channels) {
		return []float64{}, fmt.Errorf("only %v channels available, but looking for channel number %v", w.Channels, channel + 1)
	}

	for i := 0; i < buckets; i++ {
		for j := 0; j < samplesInBuckets; j++ {
			if w.BitsPerSample == uint16(8) {
				intVal, ok := w.Data[i * samplesInBuckets + j].ChannelData[channel].(uint8)
				if !ok {
					return []float64{}, fmt.Errorf("could not convert %v to int", w.Data[i * samplesInBuckets + j].ChannelData[channel])
				}
				bucketVals[i] += float64(intVal)
			}
			if w.BitsPerSample == uint16(16) {
				intVal, ok := w.Data[i * samplesInBuckets + j].ChannelData[channel].(int16)
				if !ok {
					return []float64{}, fmt.Errorf("could not convert %v to int", w.Data[i * samplesInBuckets + j].ChannelData[channel])
				}
				bucketVals[i] += float64(intVal)
			}
			if w.BitsPerSample == uint16(32) {
				intVal, ok := w.Data[i * samplesInBuckets + j].ChannelData[channel].(int32)
				if !ok {
					return []float64{}, fmt.Errorf("could not convert %v to int", w.Data[i * samplesInBuckets + j].ChannelData[channel])
				}
				bucketVals[i] += float64(intVal)
			}
			if w.BitsPerSample == uint16(64) {
				intVal, ok := w.Data[i * samplesInBuckets + j].ChannelData[channel].(int64)
				if !ok {
					return []float64{}, fmt.Errorf("could not convert %v to int", w.Data[i * samplesInBuckets + j].ChannelData[channel])
				}
				bucketVals[i] += float64(intVal)
			}
		}

		bucketVals[i] = bucketVals[i] / float64(samplesInBuckets)
	}

	return bucketVals, nil
}

// GenerateSplinePoints generates a list of points that can be used to render
// a waveform of the audio using motion canvas
func (w *Wav) GenerateSplinePoints(height, width float64, buckets, channel int) ([][]float64, error) {
	bucketedVals, err := w.GenerateBucketedAvgSampleVals(buckets, channel)
	if err != nil {
		return [][]float64{}, nil
	}

	maxAbsAmplitude, err := util.MaxAbsValue(bucketedVals)
	if err != nil {
		return [][]float64{}, nil
	}

	splinePoints := make([][]float64, buckets)
	for i, v := range bucketedVals {
		splinePoints[i] = []float64{
			float64(i) * (width / (float64(buckets) - 1)) - width / 2,
			v * (height / maxAbsAmplitude),
		}
	}

	return splinePoints, nil
}