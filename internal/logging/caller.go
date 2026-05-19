package logging

import (
	"fmt"
	"path"
	"runtime"
	"strings"
)

const (
	maxCallerDepth = 25
	minCallerDepth = 4
	loggingPackage = "logging"
)

func getPackageName(funcName string) string {
	lastPeriod := strings.LastIndex(funcName, ".")
	lastSlash := strings.LastIndex(funcName, "/")

	if lastPeriod > lastSlash {
		return funcName[:lastPeriod]
	}

	return funcName
}

func getCaller(minDepth int) *runtime.Frame {
	pcs := make([]uintptr, maxCallerDepth)
	depth := runtime.Callers(minDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for {
		frame, more := frames.Next()
		if !more {
			return nil
		}

		pkg := getPackageName(frame.Function)
		if pkg != loggingPackage {
			return &frame
		}
	}
}

func prettierFunc(frame *runtime.Frame) (function, file string) {
	file = fmt.Sprintf("%s:%v", path.Base(frame.File), frame.Line)
	function = frame.Function
	return
}
