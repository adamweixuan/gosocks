//go:build linux

package main

import (
	"golang.org/x/sys/unix"
	"os"
	"syscall"
)

func setSocketOptInt(fd int, level int, opt int, value int) {
	if err := syscall.SetsockoptInt(fd, level, opt, value); err != nil {
		log.Error("setSocketOptInt fail: %+v", os.NewSyscallError("setsockopt", err))
	}
}

func control(opt *CtrlOpt) func(network, address string, c syscall.RawConn) error { //nolint:typecheck
	return func(network, address string, c syscall.RawConn) (err error) {
		return c.Control(func(fd uintptr) {

			if opt.iface != nil {
				unix.BindToDevice(int(fd), opt.iface.Name)
			}

			if opt.reuseAddr {
				unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			}

			if opt.tcpNodelay {
				setSocketOptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
			}
		})
	}
}
