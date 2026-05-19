package utils

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func ValidateBool(value string) error {
	_, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}

	return nil
}

func ValidateTimestamp(dateTime string) error {
	_, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUInt64(value string) error {
	if value != "" {
		if _, err := strconv.ParseUint(value, 10, 64); err != nil {
			return err
		}
	}
	return nil
}

func ValidateSliceUInt64(value string) error {
	slice := strings.SplitSeq(value, ",")

	for number := range slice {
		_, err := strconv.ParseUint(number, 10, 64)
		if err != nil {
			return err
		}
	}

	return nil
}

func ValidateSemesterNumber(value string) error {
	if value != "" {
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		if v != 1 && v != 2 {
			return errors.New("semester number must be 1 or 2")
		}
	}
	return nil
}

func ValidateYear(value string) error {
	_, err := time.Parse("2006", value)
	return err
}
