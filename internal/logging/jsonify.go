package logging

import (
	"github.com/Lapakin/edu-planner/internal/adapter/json"
)

type jsonifier struct {
	Formatter FormatterType
}

func newJsonifier(formatter FormatterType) *jsonifier {
	return &jsonifier{Formatter: formatter}
}

func (j *jsonifier) jsonifyArguments(args ...any) []any {
	jsonifiedArgs := make([]any, len(args))

	for i, arg := range args {
		jsonifiedArgs[i] = j.jsonifyArgument(arg)
	}

	return jsonifiedArgs
}

func (j *jsonifier) jsonifyArgument(arg any) any {
	switch v := arg.(type) {
	case string, int8, int16, int, int32, int64, uint8, uint16, uint, uint32, uint64, float32, float64, bool, error:
		return v
	default:
		return j.jsonifyNonPrimitive(v)
	}
}

func (j *jsonifier) jsonifyNonPrimitive(v any) any {
	if j.Formatter != PrettyFormatter {
		bytes, err := json.Marshal(v)
		if err != nil {
			return v
		}
		return string(bytes)
	}

	bytes, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return v
	}
	return string(bytes)
}
