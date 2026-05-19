package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

func NewEncoder(w io.Writer) *jsoniter.Encoder {
	return jsoniter.NewEncoder(w)
}
