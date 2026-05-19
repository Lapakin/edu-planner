package utils

import (
	"strconv"
)

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func SliceUint64ToString(slice []uint64) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		result[i] = Uint64ToString(v)
	}
	return result
}
