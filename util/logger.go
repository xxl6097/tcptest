package vlog

import (
	"context"
	"fmt"
)

// 用户自定义日志方式，
// 可以通过自身业务的日志方式，来重置zinx内部引擎的日志打印方式
// 本例以fmt.Println为例
type MyLogger struct {
	OnLogRecv func(Level, string)
}

// 没有context的日志接口
func (l *MyLogger) InfoF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	if l.OnLogRecv != nil {
		l.OnLogRecv(INFO, fmt.Sprintf(format, v...))
	}
}

func (l *MyLogger) ErrorF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	if l.OnLogRecv != nil {
		l.OnLogRecv(ERROR, fmt.Sprintf(format, v...))
	}
}

func (l *MyLogger) DebugF(format string, v ...interface{}) {
	fmt.Printf(format, v...)
	if l.OnLogRecv != nil {
		l.OnLogRecv(DEBUG, fmt.Sprintf(format, v...))
	}
}

// 携带context的日志接口
func (l *MyLogger) InfoFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(format, v...)

	if l.OnLogRecv != nil {
		l.OnLogRecv(INFO, fmt.Sprintf(format, v...))
	}
}

func (l *MyLogger) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(format, v...)

	if l.OnLogRecv != nil {
		l.OnLogRecv(ERROR, fmt.Sprintf(format, v...))
	}
}

func (l *MyLogger) DebugFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(format, v...)

	if l.OnLogRecv != nil {
		l.OnLogRecv(DEBUG, fmt.Sprintf(format, v...))
	}
}
