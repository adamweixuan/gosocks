package main

import (
	"errors"
)

var (
	ErrUnSupportVersion  = errors.New("unsupport socks version")
	ErrUnsupportAddrType = errors.New("unsupport address type")
	ErrUnsupportNetType  = errors.New("unsupport net type")
	ErrNilListener       = errors.New("listener is nil")
)
