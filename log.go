package main

import (
	"context"
	stdlog "log"
	"os"

	"github.com/fatih/color"
)

type Log interface {
	CtxInfo(ctx context.Context, format string, v ...any)
	CtxWarn(ctx context.Context, format string, v ...any)
	CtxError(ctx context.Context, format string, v ...any)
	CtxTrace(ctx context.Context, format string, v ...any)

	Info(format string, v ...any)
	Warn(format string, v ...any)
	Error(format string, v ...any)
	Trace(format string, v ...any)
}

var (
	blue    = color.New(color.FgBlue).SprintFunc()    //nolint:gochecknoglobals
	red     = color.New(color.FgRed).SprintFunc()     //nolint:gochecknoglobals
	magenta = color.New(color.FgMagenta).SprintFunc() //nolint:gochecknoglobals
	green   = color.New(color.FgGreen).SprintFunc()   //nolint:gochecknoglobals

	logflag     = stdlog.LstdFlags | stdlog.Lmicroseconds       //nolint:gochecknoglobals
	infoLogger  = stdlog.New(os.Stderr, blue("INFO "), logflag) //nolint:gochecknoglobals
	warnLogger  = stdlog.New(os.Stderr, magenta("WARN "), logflag)
	errorLogger = stdlog.New(os.Stderr, red("ERROR "), logflag)
	traceLogger = stdlog.New(os.Stderr, green("TRACE "), logflag)
)

var (
	log Log = &defLogger{}
)

type defLogger struct{}

func (d *defLogger) Info(format string, v ...any) {
	infoLogger.Printf(format, v...)
}

func (d *defLogger) Warn(format string, v ...any) {
	warnLogger.Printf(format, v...)
}

func (d *defLogger) Error(format string, v ...any) {
	errorLogger.Printf(format, v...)
}

func (d *defLogger) Trace(format string, v ...any) {
	traceLogger.Printf(format, v...)
}

func (d *defLogger) CtxInfo(ctx context.Context, format string, v ...any) {
	logID, _ := ctx.Value(KeyTraceID).(string)
	infoLogger.Printf(logID+" "+format, v...)
}

func (d *defLogger) CtxWarn(ctx context.Context, format string, v ...any) {
	logID, _ := ctx.Value(KeyTraceID).(string)
	warnLogger.Printf(logID+" "+format, v...)
}

func (d *defLogger) CtxError(ctx context.Context, format string, v ...any) {
	logID, _ := ctx.Value(KeyTraceID).(string)
	errorLogger.Printf(logID+" "+format, v...)
}

func (d *defLogger) CtxTrace(ctx context.Context, format string, v ...any) {
	logID, _ := ctx.Value(KeyTraceID).(string)
	traceLogger.Printf(logID+" "+format, v...)
}
