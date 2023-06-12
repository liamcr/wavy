package util

import (
	"errors"
	"math"
)

// MaxVal finds the maximum value found in an array of float64s
func MaxVal(slice []float64) (float64, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot find max val of empty slice")
	}
	maxVal := slice[0]
	for _, v := range slice {
		maxVal = math.Max(maxVal, v)
	}

	return maxVal, nil
}

// MinVal returns the minumum value found in an array of float64s
func MinVal(slice []float64) (float64, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot find min val of empty slice")
	}
	minVal := slice[0]
	for _, v := range slice {
		minVal = math.Min(minVal, v)
	}

	return minVal, nil
}

// MaxAbsValue returns the maximum absolute value found in an array of float64s
func MaxAbsValue(slice []float64) (float64, error) {
	if len(slice) == 0 {
		return -1, errors.New("cannot find max abs val of empty slice")
	}
	maxVal := math.Abs(slice[0])
	for _, v := range slice {
		maxVal = math.Max(maxVal, math.Abs(v))
	}

	return maxVal, nil
}