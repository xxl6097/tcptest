package vlog

import (
	"fmt"
	"github.com/gookit/color"
	"path"
	"runtime"
	"strings"
	"time"
)

type Level int

var ISDEBUG = true
var LogChainIntercept func(level Level, str string)

const (
	INFO Level = iota
	DEBUG
	ERROR
	WARM
	VERBOSE
)

func getFileLine(index int, a string) string {
	_, file, lineNo, _ := runtime.Caller(index)
	fileName := path.Base(file)
	str := fmt.Sprintf("%s %s:%d: %s", time.Now().Format("2006-01-02 15:04:05.000"), fileName, lineNo, a)
	return str
}

func printf(level Level, index int, format string, v ...any) {
	if ISDEBUG {
		line := fmt.Sprintf(format, v...)
		if !strings.Contains(format, "\r") {
			line = getFileLine(index, line)
		}
		switch level {
		case INFO:
			color.HiWhite.Printf(line)
		case ERROR:
			color.HiRed.Printf(line)
		case VERBOSE:
			color.HiBlue.Printf(line)
		case DEBUG:
			color.HiGreen.Printf(line)
		case WARM:
			color.HiMagenta.Printf(line)
		default:
			color.HiWhite.Printf(line)
		}
		if LogChainIntercept != nil {
			LogChainIntercept(level, line)
		}
	}
}

func Printf(level Level, format string, v ...any) {
	printf(level, 3, format, v...)
}

func Println(level Level, v ...any) {
	printf(level, 3, "%s\n", v...)
}

func Df(format string, v ...any) {
	printf(DEBUG, 3, format, v...)
}
func Vf(format string, v ...any) {
	printf(VERBOSE, 3, format, v...)
}
func Ef(format string, v ...any) {
	printf(ERROR, 3, format, v...)
}
func If(format string, v ...any) {
	printf(INFO, 3, format, v...)
}
func Wf(format string, v ...any) {
	printf(WARM, 3, format, v...)
}

func W(v ...any) {
	printf(WARM, 3, "%s\n", v...)
}
func I(v ...any) {
	printf(INFO, 3, "%s\n", v...)
}
func E(v ...any) {
	printf(ERROR, 3, "%s\n", v...)
}
func D(v ...any) {
	printf(DEBUG, 3, "%s\n", v...)
}
func V(v ...any) {
	printf(VERBOSE, 3, "%s\n", v...)
}
