package util

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Read bytes is a helper function that will read `numBytes` from the `input`,
// and return the value of those bytes.
func ReadBytes(input io.Reader, numBytes int) ([]byte, error) {
	output := make([]byte, numBytes)
	bytesRead, err := input.Read(output)
	if err != nil {
		return []byte{}, err
	}
	if bytesRead != numBytes {
		return []byte{}, fmt.Errorf("expected %v bytes, read %v", numBytes, bytesRead)
	}

	return output, nil
}

// UInt64ToBytes takes a uint64 and returns the little endian encoded byte sequence
// representation of it
func UInt64ToBytes(num uint64) []byte {
	byteRepresentation := make([]byte, 8)
	binary.LittleEndian.PutUint64(byteRepresentation, num)
	return byteRepresentation
}

// BytesToUInt32 will take a byte array and convert it to uint32 assuming little
// endian encoding
func BytesToUInt32(bytes []byte) uint32 {
	return binary.LittleEndian.Uint32(bytes)
}

// UInt32ToBytes takes a uint32 and returns the little endian encoded byte sequence
// representation of it
func UInt32ToBytes(num uint32) []byte {
	byteRepresentation := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteRepresentation, num)
	return byteRepresentation
}

// BytesToUInt16 will take a byte array and convert it to uint16 assuming little
// endian encoding
func BytesToUInt16(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}

// UInt16ToBytes takes a uint16 and returns the little endian encoded byte sequence
// representation of it
func UInt16ToBytes(num uint16) []byte {
	byteRepresentation := make([]byte, 2)
	binary.LittleEndian.PutUint16(byteRepresentation, num)
	return byteRepresentation
}

// ReadSample reads in one sample of audio data
// For now I will only support 8, 16, 32, and 64 bit depths
// Can return either uint8, int16, int32, or int64
func ReadSample(input io.Reader, bitsPerSample uint16) (any, error) {
	bytesPerSample := int(bitsPerSample) / 8

	byteData, err := ReadBytes(input, bytesPerSample)
	if err != nil {
		return 0, err
	}

	if bitsPerSample == uint16(8) {
		// If 8 bits per sample, only one byte should have been read
		return uint8(byteData[0]), nil
	}
	if bitsPerSample == uint16(16) {
		intRepr := binary.LittleEndian.Uint16(byteData)
		return int16(intRepr), nil
	}
	if bitsPerSample == uint16(32) {
		intRepr := binary.LittleEndian.Uint32(byteData)
		return int32(intRepr), nil
	}
	if bitsPerSample == uint16(64) {
		intRepr := binary.LittleEndian.Uint64(byteData)
		return int64(intRepr), nil
	}

	return 0, fmt.Errorf("bit depth not one of 8, 16, 32, or 64 (%d)", bitsPerSample)
}