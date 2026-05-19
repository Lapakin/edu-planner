package json

import (
	"io"

	jsoniter "github.com/json-iterator/go"
)

func NewDecoder(r io.Reader) *jsoniter.Decoder {
	return jsoniter.NewDecoder(r)
}
