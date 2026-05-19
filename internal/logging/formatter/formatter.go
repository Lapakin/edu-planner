package formatter

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mgutz/ansi"
	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

const defaultTimestampFormat = time.RFC3339

var (
	defaultColorScheme = &ColorScheme{
		InfoLevelStyle:  "white+bh:green",
		WarnLevelStyle:  "white+bh:yellow",
		ErrorLevelStyle: "white+bh:red",
		FatalLevelStyle: "white+bh:red",
		PanicLevelStyle: "white+bh:red",
		DebugLevelStyle: "white+bh:blue",
		TraceLevelStyle: "white+bh:magenta",
		TimestampStyle:  "black+h",
	}
	noColorsColorScheme = &compiledColorScheme{
		InfoLevelColor:  ansi.ColorFunc(""),
		WarnLevelColor:  ansi.ColorFunc(""),
		ErrorLevelColor: ansi.ColorFunc(""),
		FatalLevelColor: ansi.ColorFunc(""),
		PanicLevelColor: ansi.ColorFunc(""),
		DebugLevelColor: ansi.ColorFunc(""),
		TraceLevelStyle: ansi.ColorFunc(""),
		TimestampColor:  ansi.ColorFunc(""),
	}
	defaultCompiledColorScheme = compileColorScheme(defaultColorScheme)
)

type ColorScheme struct {
	InfoLevelStyle  string
	WarnLevelStyle  string
	ErrorLevelStyle string
	FatalLevelStyle string
	PanicLevelStyle string
	DebugLevelStyle string
	TraceLevelStyle string
	TimestampStyle  string
}

type compiledColorScheme struct {
	InfoLevelColor  func(string) string
	WarnLevelColor  func(string) string
	ErrorLevelColor func(string) string
	FatalLevelColor func(string) string
	PanicLevelColor func(string) string
	DebugLevelColor func(string) string
	TraceLevelStyle func(string) string
	TimestampColor  func(string) string
}

func compileColorScheme(s *ColorScheme) *compiledColorScheme {
	return &compiledColorScheme{
		InfoLevelColor:  getCompiledColor(s.InfoLevelStyle, defaultColorScheme.InfoLevelStyle),
		WarnLevelColor:  getCompiledColor(s.WarnLevelStyle, defaultColorScheme.WarnLevelStyle),
		ErrorLevelColor: getCompiledColor(s.ErrorLevelStyle, defaultColorScheme.ErrorLevelStyle),
		FatalLevelColor: getCompiledColor(s.FatalLevelStyle, defaultColorScheme.FatalLevelStyle),
		PanicLevelColor: getCompiledColor(s.PanicLevelStyle, defaultColorScheme.PanicLevelStyle),
		DebugLevelColor: getCompiledColor(s.DebugLevelStyle, defaultColorScheme.DebugLevelStyle),
		TraceLevelStyle: getCompiledColor(s.TraceLevelStyle, defaultColorScheme.TraceLevelStyle),
		TimestampColor:  getCompiledColor(s.TimestampStyle, defaultColorScheme.TimestampStyle),
	}
}

func getCompiledColor(main, fallback string) func(string) string {
	if main != "" {
		return ansi.ColorFunc(main)
	}

	return ansi.ColorFunc(fallback)
}

type PrettyFormatter struct {
	ForceColors      bool
	DisableColors    bool
	ForceFormatting  bool
	DisableTimestamp bool
	DisableUppercase bool
	FullTimestamp    bool
	TimestampFormat  string
	DisableSorting   bool
	SpacePadding     int
	colorScheme      *compiledColorScheme
	isTerminal       bool
	CallerPrettyfier func(*runtime.Frame) (function string, file string)
	once             sync.Once
}

func (f *PrettyFormatter) init(entry *logrus.Entry) {
	if entry.Logger != nil {
		f.isTerminal = f.checkIfTerminal(entry.Logger.Out)
	}
}

func (f *PrettyFormatter) checkIfTerminal(w io.Writer) bool {
	switch v := w.(type) {
	case *os.File:
		return term.IsTerminal(int(v.Fd()))
	default:
		return false
	}
}

func (f *PrettyFormatter) SetColorScheme(colorScheme *ColorScheme) {
	f.colorScheme = compileColorScheme(colorScheme)
}

func (f *PrettyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	if entry.Buffer != nil {
		b = entry.Buffer
	}

	f.once.Do(func() { f.init(entry) })

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if entry.HasCaller() {
		funcVal, fileVal := entry.Caller.Function, fmt.Sprintf("%s:%d", entry.Caller.File, entry.Caller.Line)
		if f.CallerPrettyfier != nil {
			funcVal, fileVal = f.CallerPrettyfier(entry.Caller)
		}

		if funcVal != "" {
			entry.Data["func"] = funcVal
		}
		if fileVal != "" {
			entry.Data["file"] = fileVal
		}
	}

	isColored := (f.ForceColors || f.isTerminal) && !f.DisableColors
	colorScheme := noColorsColorScheme
	if isColored {
		colorScheme = defaultCompiledColorScheme

		if f.colorScheme != nil {
			colorScheme = f.colorScheme
		}
	}
	f.printColored(b, entry, timestampFormat, colorScheme)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *PrettyFormatter) printColored(b *bytes.Buffer, entry *logrus.Entry, timestampFormat string, colorScheme *compiledColorScheme) {
	levelColor := getLevelColor(entry.Level, colorScheme)
	levelText := getLevelText(entry.Level, f.DisableUppercase)
	level := levelColor(fmt.Sprintf(" %s ", levelText))
	message := entry.Message

	messageFormat := "%s"
	if f.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", f.SpacePadding)
	}

	app, postfix := getAppAndPostfix(entry.Data, colorScheme.TimestampColor)

	if f.DisableTimestamp {
		format := "%s %s: " + messageFormat + "%s"
		_, _ = fmt.Fprintf(b, format, app, level, message, postfix)
		return
	}

	timestamp := fmt.Sprintf("[%s]", entry.Time.Format(timestampFormat))
	format := "%s %s %s: " + messageFormat + " %s"
	_, _ = fmt.Fprintf(b, format, colorScheme.TimestampColor(timestamp), app, level, message, postfix)
}

func getLevelColor(level logrus.Level, colorScheme *compiledColorScheme) func(string) string {
	switch level {
	case logrus.TraceLevel:
		return colorScheme.TraceLevelStyle
	case logrus.DebugLevel:
		return colorScheme.DebugLevelColor
	case logrus.InfoLevel:
		return colorScheme.InfoLevelColor
	case logrus.WarnLevel:
		return colorScheme.WarnLevelColor
	case logrus.ErrorLevel:
		return colorScheme.ErrorLevelColor
	case logrus.FatalLevel:
		return colorScheme.FatalLevelColor
	case logrus.PanicLevel:
		return colorScheme.PanicLevelColor
	default:
		return colorScheme.DebugLevelColor
	}
}

func getLevelText(level logrus.Level, disableUppercase bool) string {
	levelText := level.String()
	if !disableUppercase {
		levelText = strings.ToUpper(levelText)
	}
	return levelText
}

func getAppAndPostfix(fields logrus.Fields, timestampColor func(string) string) (app, postfix string) {
	for k, v := range fields {
		if k == "app" {
			app = timestampColor(fmt.Sprintf("[%s]", v))
			continue
		}

		if postfix != "" {
			postfix += " "
		}
		postfix += fmt.Sprintf("%s=%+v", k, v)
	}
	postfix = timestampColor(fmt.Sprintf("(%s)", postfix))
	return app, postfix
}
