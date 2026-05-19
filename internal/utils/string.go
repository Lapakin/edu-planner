package utils

import (
	"strconv"
)

func StringToUint64(value string) (uint64, error) {
	number, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return number, err
	}
	return number, nil
}
