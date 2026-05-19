package logging

func (l *Logger) withCallerFields(minDepth int) *Logger {
	function, file := prettierFunc(getCaller(minDepth))
	return l.WithField("func", function).
		WithField("file", file)
}

func (l *Logger) Trace(args ...any) {
	l.withCallerFields(minCallerDepth).Log(TraceLevel, args...)
}

func (l *Logger) Debug(args ...any) {
	l.withCallerFields(minCallerDepth).Log(DebugLevel, args...)
}

func (l *Logger) Print(args ...any) {
	l.withCallerFields(minCallerDepth).Info(args...)
}

func (l *Logger) Info(args ...any) {
	l.withCallerFields(minCallerDepth).Log(InfoLevel, args...)
}

func (l *Logger) Warn(args ...any) {
	l.withCallerFields(minCallerDepth).Log(WarnLevel, args...)
}

func (l *Logger) Warning(args ...any) {
	l.withCallerFields(minCallerDepth + 1).Warn(args...)
}

func (l *Logger) Error(args ...any) {
	l.withCallerFields(minCallerDepth).Log(ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...any) {
	l.withCallerFields(minCallerDepth).Log(FatalLevel, args...)
	l.Logger.Exit(1)
}

func (l *Logger) Panic(args ...any) {
	l.withCallerFields(minCallerDepth).Log(PanicLevel, args...)
}

func (l *Logger) Tracef(format string, args ...any) {
	l.withCallerFields(minCallerDepth).Logf(TraceLevel, format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.withCallerFields(minCallerDepth).Logf(DebugLevel, format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Printf(format string, args ...any) {
	l.withCallerFields(minCallerDepth+1).Infof(format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Infof(format string, args ...any) {
	l.withCallerFields(minCallerDepth).Logf(InfoLevel, format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Warnf(format string, args ...any) {
	l.withCallerFields(minCallerDepth).Logf(WarnLevel, format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Warningf(format string, args ...any) {
	l.withCallerFields(minCallerDepth+1).Warnf(format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.withCallerFields(minCallerDepth).Logf(ErrorLevel, format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.withCallerFields(minCallerDepth).Logf(FatalLevel, format, l.jsonifyArguments(args...)...)
	l.Logger.Exit(1)
}

func (l *Logger) Panicf(format string, args ...any) {
	l.withCallerFields(minCallerDepth).Logf(PanicLevel, format, l.jsonifyArguments(args...)...)
}

func (l *Logger) Traceln(args ...any) {
	l.withCallerFields(minCallerDepth).Logln(TraceLevel, args...)
}

func (l *Logger) Debugln(args ...any) {
	l.withCallerFields(minCallerDepth).Logln(DebugLevel, args...)
}

func (l *Logger) Infoln(args ...any) {
	l.withCallerFields(minCallerDepth).Logln(InfoLevel, args...)
}

func (l *Logger) Println(args ...any) {
	l.withCallerFields(minCallerDepth + 1).Infoln(args...)
}

func (l *Logger) Warnln(args ...any) {
	l.withCallerFields(minCallerDepth).Logln(WarnLevel, args...)
}

func (l *Logger) Warningln(args ...any) {
	l.withCallerFields(minCallerDepth + 1).Warnln(args...)
}

func (l *Logger) Errorln(args ...any) {
	l.withCallerFields(minCallerDepth).Logln(ErrorLevel, args...)
}

func (l *Logger) Fatalln(args ...any) {
	l.withCallerFields(minCallerDepth).Logln(FatalLevel, args...)
	l.Logger.Exit(1)
}

func (l *Logger) Panicln(args ...any) {
	l.withCallerFields(minCallerDepth).Logln(PanicLevel, args...)
}
