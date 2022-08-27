//go:build !linux && !darwin

package main

import (
	"syscall"
)

func control(opt *CtrlOpt) func(string, string, syscall.RawConn) error { return nil }
