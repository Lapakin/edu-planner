package logging

import (
	"github.com/sirupsen/logrus"

	"github.com/Lapakin/edu-planner/internal/logging/formatter"
)

type Logger struct {
	*logrus.Entry
	*jsonifier
}

func NewLogger(level Level, formatterType FormatterType) *Logger {
	l := logrus.New()
	l.SetLevel(level)

	switch formatterType {
	case JSONFormatter:
		l.Formatter = &logrus.JSONFormatter{}
	case PrettyFormatter:
		l.Formatter = &formatter.PrettyFormatter{
			ForceColors:      false,
			DisableColors:    false,
			ForceFormatting:  false,
			DisableTimestamp: false,
			DisableUppercase: false,
			FullTimestamp:    true,
			TimestampFormat:  "2006-01-02 15:04:05",
			DisableSorting:   false,
			SpacePadding:     0,
			CallerPrettyfier: nil,
		}
	default:
		l.Formatter = &logrus.TextFormatter{
			FullTimestamp: true,
			DisableColors: true,
		}
	}

	return &Logger{
		Entry:     logrus.NewEntry(l),
		jsonifier: newJsonifier(formatterType),
	}
}

func (l *Logger) WithField(field string, value any) *Logger {
	return &Logger{
		Entry:     l.Entry.WithField(field, value),
		jsonifier: l.jsonifier,
	}
}

type FormatterType string

const (
	JSONFormatter   = "JSON"
	TextFormatter   = "TextFormatter"
	PrettyFormatter = "PrettyFormatter"
)
