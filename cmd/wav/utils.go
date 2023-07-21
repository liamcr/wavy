package wav

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/liamcr/wavy/internal/util"
)

// GenerateBucketedAvgSampleVals separates the audio file data into a set number of
// buckets, by dividing the data array into equal parts, and averaging the amplitudes
// over each part
func (w *Wav) GenerateBucketedAvgSampleVals(buckets, channel int, abs bool) ([]float64, error) {
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

				if abs {
					bucketVals[i] += math.Abs(float64(intVal))
				} else {
					bucketVals[i] += float64(intVal)
				}
			}
			if w.BitsPerSample == uint16(16) {
				intVal, ok := w.Data[i * samplesInBuckets + j].ChannelData[channel].(int16)
				if !ok {
					return []float64{}, fmt.Errorf("could not convert %v to int", w.Data[i * samplesInBuckets + j].ChannelData[channel])
				}

				if abs {
					bucketVals[i] += math.Abs(float64(intVal))
				} else {
					bucketVals[i] += float64(intVal)
				}
			}
			if w.BitsPerSample == uint16(32) {
				intVal, ok := w.Data[i * samplesInBuckets + j].ChannelData[channel].(int32)
				if !ok {
					return []float64{}, fmt.Errorf("could not convert %v to int", w.Data[i * samplesInBuckets + j].ChannelData[channel])
				}
				
				if abs {
					bucketVals[i] += math.Abs(float64(intVal))
				} else {
					bucketVals[i] += float64(intVal)
				}
			}
			if w.BitsPerSample == uint16(64) {
				intVal, ok := w.Data[i * samplesInBuckets + j].ChannelData[channel].(int64)
				if !ok {
					return []float64{}, fmt.Errorf("could not convert %v to int", w.Data[i * samplesInBuckets + j].ChannelData[channel])
				}

				if abs {
					bucketVals[i] += math.Abs(float64(intVal))
				} else {
					bucketVals[i] += float64(intVal)
				}
			}
		}

		bucketVals[i] = bucketVals[i] / float64(samplesInBuckets)
	}

	return bucketVals, nil
}

// GenerateSplinePoints generates a list of points that can be used to render
// a waveform of the audio using motion canvas
func (w *Wav) GenerateSplinePoints(height, width float64, buckets, channel int) ([][]float64, error) {
	bucketedVals, err := w.GenerateBucketedAvgSampleVals(buckets, channel, false)
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

// CastToInt will a variable input value which could be uint8, int16, int32
// or int64 and normalize it to be an int
func CastToInt(v any) (int, error) {
	uint8Val, ok := v.(uint8)
	if ok {
		return int(uint8Val), nil
	}
	int16Val, ok := v.(int16)
	if ok {
		return int(int16Val), nil
	}
	int32Val, ok := v.(int32)
	if ok {
		return int(int32Val), nil
	}
	int64Val, ok := v.(int64)
	if ok {
		return int(int64Val), nil
	}

	return 0, errors.New("cannot convert value to int")
		
}

// Const vals representing GenerateSvg config
const pathTemplate = "<rect x=\"%d\" y=\"%d\" width=\"%d\" height=\"%d\" fill=\"%s\" ry=\"%d\" rx=\"%d\"/>"
const svgHeight = 100
const desiredPillWidth = 12
const desiredPillMargin = 2
const desiredSVGPadding = 2

func (w *Wav) GenerateSvg(outputPath string, width, channel int, fill string) (*os.File, error) {
	numBuckets := (width - 2 * desiredSVGPadding) / (desiredPillWidth + desiredPillMargin)
	bucketedVals, err := w.GenerateBucketedAvgSampleVals(numBuckets, channel, false)
	if err != nil {
		return nil, err
	}

	fmt.Println(bucketedVals[0])

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}

	if _, err := outputFile.Write([]byte(fmt.Sprintf("<svg xmlns=\"http://www.w3.org/2000/svg\" width=\"100%%\" height=\"100%%\" viewBox=\"0 -50 %d %d\">", width, svgHeight))); err != nil {
		return nil, err
	}

	maxHeight, err := util.MaxVal(bucketedVals)
	if err != nil {
		return nil, err
	}

	for i, bucketVal := range bucketedVals {
		pillHeight := (bucketVal / maxHeight) * float64((svgHeight - desiredSVGPadding))
		if _, err := outputFile.Write([]byte(
			fmt.Sprintf(
				pathTemplate,
				desiredSVGPadding + (i * (desiredPillMargin + desiredPillWidth)),
				int(math.Min((-1 * (pillHeight / 2)), float64(-1 * desiredPillWidth / 2))),
				desiredPillWidth,
				int(math.Max(pillHeight, float64(desiredPillWidth))),
				fill,
				desiredPillWidth / 2,
				desiredPillWidth / 2,
			),
		)); err != nil {
			return nil, err
		}
	}

	if _, err := outputFile.Write([]byte("</svg>")); err != nil {
		return nil, err
	}

	return outputFile, nil
}