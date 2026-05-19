package json

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
