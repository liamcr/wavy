package wav

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/liamcr/wavy/internal/util"
)

const (
	fmtSize = 16
	chunkHeadingSize = 8
)

// SampleGroup is a representation of the group of samples that represent one
// moment in time, with each sample in the group representing a channel.
type SampleGroup struct {
	// ChannelData represents the sample values taken in each
	// audio channel.
	ChannelData []any
}

// Wav is a struct representation of a .wav audio file
type Wav struct {
	// FormatType is the type of audio format (1 = PCM)
	FormatType uint16

	// Channels is the number of channels
	Channels uint16

	// SampleRate is the number of samples per second. Common values are
	// 44100 (CD) or 48000 (DAT)
	SampleRate uint32

	// DataRate is the average number of bytes per second. Equal to
	// (Sample Rate * BitsPerSample * Channels) / 8
	DataRate uint32

	// DataBlockSize is the minimum atomic unit of data, in bytes
	DataBlockSize uint16

	// BitsPerSample is the number of bits per sample (bit depth)
	BitsPerSample uint16

	// DataSize is the size in bytes of the audio data
	DataSize uint32
	
	// Data is an array of the sample data parsed from the Wav file
	Data []SampleGroup
}

// Encode will take the attributes found in the parent struct and will output
// a byte representation of a valid wav file.
func (w *Wav) Encode() ([]byte, error) {
	encodedWav := []byte("RIFF")

	fileSize := 3 * chunkHeadingSize + fmtSize + w.DataSize + 4
	encodedWav = append(encodedWav, util.UInt32ToBytes(uint32(fileSize))...)

	encodedWav = append(encodedWav, "WAVE"...)

	// fmt chunk
	encodedWav = append(encodedWav, "fmt "...)

	// fmt chunk should always be 16 bytes long
	encodedWav = append(encodedWav, util.UInt32ToBytes(uint32(16))...)
	encodedWav = append(encodedWav, util.UInt16ToBytes(w.FormatType)...)
	encodedWav = append(encodedWav, util.UInt16ToBytes(w.Channels)...)
	encodedWav = append(encodedWav, util.UInt32ToBytes(w.SampleRate)...)
	encodedWav = append(encodedWav, util.UInt32ToBytes(w.DataRate)...)
	encodedWav = append(encodedWav, util.UInt16ToBytes(w.DataBlockSize)...)
	encodedWav = append(encodedWav, util.UInt16ToBytes(w.BitsPerSample)...)

	encodedWav = append(encodedWav, "data"...)
	encodedWav = append(encodedWav, util.UInt32ToBytes(w.DataSize)...)
	
	// Write data
	for _, sampleGroup := range(w.Data) {
		for _, sample := range(sampleGroup.ChannelData) {
			if w.BitsPerSample == uint16(8) {
				eightBitSample, ok := sample.(uint8)
				if !ok {
					return nil, fmt.Errorf("can't cast data point %v to byte", sample)
				}
				encodedWav = append(encodedWav, byte(eightBitSample))
			} else if w.BitsPerSample == uint16(16) {
				sixteenBitSample, ok := sample.(int16)
				if !ok {
					return nil, fmt.Errorf("can't cast data point %v to int16", sample)
				}
				encodedWav = append(encodedWav, util.UInt16ToBytes(uint16(sixteenBitSample))...)
			} else if w.BitsPerSample == uint16(32) {
				thirtyTwoBitSample, ok := sample.(int32)
				if !ok {
					return nil, fmt.Errorf("can't cast data point %v to int32", sample)
				}
				encodedWav = append(encodedWav, util.UInt32ToBytes(uint32(thirtyTwoBitSample))...)
			} else if w.BitsPerSample == uint16(64) {
				sixtyFourBitSample, ok := sample.(int64)
				if !ok {
					return nil, fmt.Errorf("can't cast data point %v to int64", sample)
				}
				encodedWav = append(encodedWav, util.UInt64ToBytes(uint64(sixtyFourBitSample))...)
			}
		}
	}

	return encodedWav, nil
}

func (w *Wav) Write(filename string) error {
	bytes, err := w.Encode()
	if err != nil {
        return err
    }

	output, err := os.Create(filename)
    if err != nil {
        return err
    }
	defer func() {
		if err := output.Close(); err != nil {
			panic(fmt.Sprintf("closing file: %v", err.Error()))
		}
	}()

	if _, err := output.Write(bytes); err != nil {
		return err
	}

	return nil
}

// Decode will take an input wav file and return a `Wav` struct with fields representing
// each attribute of the file.
func Decode(input io.Reader) (*Wav, error) {
	decodedWav := &Wav{}
	riff, err := util.ReadBytes(input, 4)
	if err != nil {
		return nil, err
	}
	if string(riff) != "RIFF" {
		return nil, errors.New("corrupted file, first 4 bytes not 'RIFF'")
	}

	_, err = util.ReadBytes(input, 4)
	if err != nil {
		return nil, err
	}

	wave, err := util.ReadBytes(input, 4)
	if err != nil {
		return nil, err
	}
	if string(wave) != "WAVE" {
		return nil, errors.New("corrupted file, bytes 9-12 do not read 'WAVE'")
	}

	// Scan through the chunks of the file, we only care about the fmt and data
	// chunks
	for {
		chunkHeader, err := util.ReadBytes(input, 4)
		if err != nil {
			return nil, err
		}

		chunkSizeBytes, err := util.ReadBytes(input, 4)
		if err != nil {
			return nil, err
		}

		chunkSize := util.BytesToUInt32(chunkSizeBytes)

		if string(chunkHeader) == "fmt " {
			err = readFmtChunk(input, decodedWav)
			if err != nil {
				return nil, err
			}
		} else if string(chunkHeader) == "data" {
			decodedWav.DataSize = chunkSize
			err = readDataChunk(input, decodedWav, chunkSize)
			if err != nil {
				return nil, err
			}
		} else {
			// Skip to the next chunk
			_, err = util.ReadBytes(input, int(chunkSize))
			if err != nil {
				return nil, err
			}
		}

		if decodedWav.Channels != 0 && len(decodedWav.Data) > 0 {
			break
		}
	}
	
	return decodedWav, nil
}

// GetDuration gets the duration of the audio file in seconds
func (w *Wav) GetDuration() float64 {
	bytesPerSample := float64(w.BitsPerSample) / 8
	return float64(w.DataSize) / bytesPerSample / float64(w.SampleRate) / float64(w.Channels)
}

func readFmtChunk(input io.Reader, wav *Wav) error {
	formatType, err := util.ReadBytes(input, 2)
	if err != nil {
		return err
	}
	wav.FormatType = util.BytesToUInt16(formatType)

	numChannels, err := util.ReadBytes(input, 2)
	if err != nil {
		return err
	}
	wav.Channels = util.BytesToUInt16(numChannels)

	sampleRate, err := util.ReadBytes(input, 4)
	if err != nil {
		return err
	}
	wav.SampleRate = util.BytesToUInt32(sampleRate)

	dataRate, err := util.ReadBytes(input, 4)
	if err != nil {
		return err
	}
	wav.DataRate = util.BytesToUInt32(dataRate)

	dataBlockSize, err := util.ReadBytes(input, 2)
	if err != nil {
		return err
	}
	wav.DataBlockSize = util.BytesToUInt16(dataBlockSize)

	bitsPerSample, err := util.ReadBytes(input, 2)
	if err != nil {
		return err
	}
	wav.BitsPerSample = util.BytesToUInt16(bitsPerSample)

	if wav.BitsPerSample != 16 {
		return fmt.Errorf("only 16-bit wav files are currently supported (current bits/sample = %v)", wav.BitsPerSample)
	}

	return nil
}

func readDataChunk(input io.Reader, wav *Wav, chunkSize uint32) error {
	dataPoints := []SampleGroup{}
	dataSize := int(wav.DataSize)
	bytesPerSampleGroup := int(int(wav.BitsPerSample) / 8) * int(wav.Channels)
	for position := 0; position < dataSize; position += bytesPerSampleGroup {
		newDataPoint := SampleGroup{}
		for i := 0; i < int(wav.Channels); i++ {
			sample, err := util.ReadSample(input, wav.BitsPerSample)
			if err != nil {
				return err
			}

			newDataPoint.ChannelData = append(newDataPoint.ChannelData, sample)
		}

		dataPoints = append(dataPoints, newDataPoint)
	}

	wav.Data = dataPoints
	return nil
}