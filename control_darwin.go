//go:build darwin

package main

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
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
				switch network {
				case "tcp4", "udp4":
					setSocketOptInt(int(fd), unix.IPPROTO_IP, unix.IP_BOUND_IF, opt.iface.Index)
				case "tcp6", "udp6":
					setSocketOptInt(int(fd), unix.IPPROTO_IPV6, unix.IPV6_BOUND_IF, opt.iface.Index)
				}
			}
			if opt.reuseAddr {
				setSocketOptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				setSocketOptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			}

			if opt.tcpNodelay {
				setSocketOptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_NODELAY, 1)
			}

		})
	}
}
