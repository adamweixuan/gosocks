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
	blue    = color.New(color.FgBlue).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()

	logflag     = stdlog.LstdFlags | stdlog.Lmicroseconds
	infoLogger  = stdlog.New(os.Stderr, blue("INFO "), logflag)
	warnLogger  = stdlog.New(os.Stderr, magenta("WARN "), logflag)
	errorLogger = stdlog.New(os.Stderr, red("ERROR "), logflag)
	traceLogger = stdlog.New(os.Stderr, green("TRACE "), logflag)
)

var (
	log Log = &defLogger{} //nolint:typecheck
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
	logId, _ := ctx.Value(KeyTraceID).(string)
	infoLogger.Printf(logId+" "+format, v...)
}

func (d *defLogger) CtxWarn(ctx context.Context, format string, v ...any) {
	logId, _ := ctx.Value(KeyTraceID).(string)
	warnLogger.Printf(logId+" "+format, v...)
}

func (d *defLogger) CtxError(ctx context.Context, format string, v ...any) {
	logId, _ := ctx.Value(KeyTraceID).(string)
	errorLogger.Printf(logId+" "+format, v...)
}

func (d *defLogger) CtxTrace(ctx context.Context, format string, v ...any) {
	logId, _ := ctx.Value(KeyTraceID).(string)
	traceLogger.Printf(logId+" "+format, v...)
}
