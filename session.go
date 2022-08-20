package main

import (
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
)

var (
	ErrUnSupportVersion  = errors.New("unsupport socks version")
	ErrUnsupportAddrType = errors.New("unsupport address type")
)

type Addrtype uint8

const (
	_      Addrtype = iota
	ipv4            = 1
	domain          = 3
	ipv6            = 4
)

const (
	socksVersion = 0x05
	connectCmd   = 0x01
)

const (
	maxBufSize = 4096
	minBufSize = 256
)

var (
	successReply = []byte{
		0x05,
		0x00,
		0x00,
		0x01,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
		0x00,
	}

	authReply = []byte{
		0x05,
		0x00,
	}
)

type Session struct {
	verbos bool
	local  net.Conn
	remote net.Conn
	id     uint64
}

func (session *Session) Start() {
	if err := session.auth(); err != nil {
		log.Printf("session auth error: %s", err.Error())
		_ = session.local.Close()
		return
	}

	if session.verbos {
		log.Printf("auth %s success", session.local.LocalAddr().String())
	}

	if err := session.connect(); err != nil {
		log.Printf("session connect error: %s", err.Error())
		_ = session.local.Close()
		return
	}
	if session.verbos {
		log.Printf("conect to %s success", session.remote.RemoteAddr().String())
	}
	session.relay()
}

func (session *Session) auth() error {
	buf := make([]byte, minBufSize)
	if _, err := io.ReadAtLeast(session.local, buf[:2], 2); err != nil {
		return err
	}
	ver, methodCnt := int(buf[0]), int(buf[1])
	if ver != socksVersion {
		return ErrUnSupportVersion
	}
	if _, err := io.ReadAtLeast(session.local, buf[:methodCnt], methodCnt); err != nil {
		return err
	}
	_, err := session.local.Write(authReply)
	return err
}

func (session *Session) connect() error {
	buf := make([]byte, maxBufSize)

	if _, err := io.ReadAtLeast(session.local, buf[:4], 4); err != nil {
		return err
	}

	ver, cmd, atyp := buf[0], buf[1], buf[3]

	if ver != socksVersion || cmd != connectCmd {
		return ErrUnSupportVersion
	}

	addr := ""
	switch atyp {
	case ipv4:
		if _, err := io.ReadAtLeast(session.local, buf[:net.IPv4len], net.IPv4len); err != nil {
			return err
		}
		addr = net.IP(buf[:net.IPv4len]).String()
	case domain:
		if _, err := io.ReadAtLeast(session.local, buf[:1], 1); err != nil {
			return err
		}
		addrlen := int(buf[0])

		if _, err := io.ReadAtLeast(session.local, buf[:addrlen], addrlen); err != nil {
			return err
		}
		addr = string(buf[:addrlen])
	case ipv6:
		if _, err := io.ReadAtLeast(session.local, buf[:net.IPv6len], net.IPv6len); err != nil {
			return err
		}
		addr = net.IP(buf[:net.IPv6len]).String()
	default:
		return ErrUnsupportAddrType
	}

	if _, err := io.ReadAtLeast(session.local, buf[:2], 2); err != nil {
		return err
	}

	port := binary.BigEndian.Uint16(buf[:2])

	endPoint := net.JoinHostPort(addr, strconv.Itoa(int(port)))
	dst, err := net.Dial("tcp", endPoint)
	if err != nil {
		_ = dst.Close()
		return err
	}
	if _, err := session.local.Write(successReply); err != nil {
		_ = dst.Close()
		return err
	}
	if session.verbos {
		log.Printf("connect to %s success", endPoint)
	}
	session.remote = dst
	return nil
}

func (session *Session) relay() {
	go func() {
		if cnt, err := io.Copy(session.local, session.remote); err != nil {
			log.Printf("session_id:%d, relay error: %s", session.id, err.Error())
		} else if session.verbos {
			log.Printf("session_id:%d:%s<<<--%s.size:%d",
				session.id, session.remote.RemoteAddr().String(),
				session.local.LocalAddr().String(), cnt)
		}
	}()

	go func() {
		if cnt, err := io.Copy(session.remote, session.local); err != nil {
			log.Printf("session_id:%d, relay error: %s", session.id, err.Error())
		} else if session.verbos {
			log.Printf("session_id:%d:%s--->>>%s.size:%d",
				session.id, session.local.LocalAddr().String(),
				session.remote.RemoteAddr().String(), cnt)
		}
	}()
}
