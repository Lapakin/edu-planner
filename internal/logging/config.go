package logging

type Config struct {
	Level     string        `yaml:"level"`
	Formatter FormatterType `yaml:"formatter"`
}
