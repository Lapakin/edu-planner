package test

import (
	"os"
)

var (
	UnitTesting         TestingType = "unit"
	IntegrationTesting  TestingType = "integration"
	NotSpecifiedTesting TestingType = "not_specified"
)

type TestingType string

func GetTestingType() TestingType {
	v := os.Getenv("INTEGRATION")
	switch v {
	case "true":
		return IntegrationTesting
	case "false":
		return UnitTesting
	default:
		return NotSpecifiedTesting
	}
}
