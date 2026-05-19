package yaml

import (
	"os"

	"gopkg.in/yaml.v3"
)

func UnmarshalYAML(path string, dest any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, dest)
}
