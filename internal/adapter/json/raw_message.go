package json

import jsonstd "encoding/json"

type RawMessage jsonstd.RawMessage

func (m *RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil || *m == nil {
		return []byte("null"), nil
	}
	return *m, nil
}

func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if data == nil {
		*m = nil
		return nil
	}
	*m = append((*m)[0:0], data...)
	return nil
}
