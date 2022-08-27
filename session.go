package main

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"time"
)

const (
	socksVersion = 0x05
	connectCmd   = 0x01
)

const (
	maxBufSize     = 4096
	minBufSize     = 256
	defaultBufSize = 32 * 1024
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
	verbos    bool
	local     net.Conn
	remote    net.Conn
	ip        net.IP
	iface     *net.Interface
	ntype     Network
	timeout   time.Duration
	noDelay   bool
	reuseAddr bool
}

func (session *Session) Start(ctx context.Context) {
	if err := session.auth(); err != nil {
		log.CtxError(ctx, "session auth error: %s", err.Error())
		session.closeConn(ctx, session.local)
		return
	}

	if session.verbos {
		log.CtxTrace(ctx, "auth %s success", session.local.LocalAddr().String())
	}

	if err := session.connect(ctx); err != nil {
		log.CtxError(ctx, "session connect error: %s", err.Error())
		session.closeConn(ctx, session.local)
		return
	}
	if session.verbos {
		log.CtxTrace(ctx, "connect to %s success", session.remote.RemoteAddr().String())
	}
	session.relay(ctx)
}

func (session *Session) auth() error {
	buf := Get(minBufSize, minBufSize)
	defer func() {
		Put(buf)
	}()
	if _, err := io.ReadAtLeast(session.local, buf[:2], 2); err != nil {
		return err
	}
	ver, mCnt := int(buf[0]), int(buf[1])
	if ver != socksVersion {
		return ErrUnSupportVersion
	}
	if _, err := io.ReadAtLeast(session.local, buf[:mCnt], mCnt); err != nil {
		return err
	}
	_, err := session.local.Write(authReply)
	return err
}

func (session *Session) connect(ctx context.Context) error {
	buf := Get(maxBufSize, maxBufSize)
	defer func() {
		Put(buf)
	}()

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
	dst, err := session.dial(endPoint)
	if err != nil {
		session.closeConn(ctx, dst)
		return err
	}
	if _, err := session.local.Write(successReply); err != nil {
		session.closeConn(ctx, dst)
		return err
	}
	session.remote = dst
	return nil
}

func (session *Session) relay(ctx context.Context) {
	forward := func(from, to net.Conn) {
		buf := Get(defaultBufSize, defaultBufSize)
		defer Put(buf)
		//cnt, err := io.Copy(from, to)
		cnt, err := io.CopyBuffer(from, to, buf)
		if err != nil {
			log.CtxError(ctx, "[%s->%s] forward error: %s", from.RemoteAddr(), to.RemoteAddr(), err.Error())
		}
		if session.verbos {
			log.CtxTrace(ctx, "[%s->%s] forward success. copy size:%d ", from.RemoteAddr(), to.RemoteAddr(), cnt)
		}

	}
	go forward(session.local, session.remote)
	go forward(session.remote, session.local)
}

func (session *Session) dial(addr string) (net.Conn, error) {
	var la net.Addr
	switch session.ntype {
	case tcp:
		la = &net.TCPAddr{
			IP: session.ip,
		}
	case udp:
		la = &net.UDPAddr{
			IP: session.ip,
		}
	default:
		return nil, ErrUnsupportNetType
	}

	dialer := &net.Dialer{LocalAddr: la, Timeout: session.timeout}

	var ctrlOps []CtrlOp

	if session.iface != nil {
		ctrlOps = append(ctrlOps, Bind(session.iface))
	}

	if session.noDelay {
		ctrlOps = append(ctrlOps, NoDelay())
	}

	if session.reuseAddr {
		ctrlOps = append(ctrlOps, ReuseAddr())
	}
	dialer.Control = Control(ctrlOps...)

	conn, err := dialer.Dial(session.ntype.String(), addr)

	if conn == nil || err != nil {
		return nil, err
	}

	if c, ok := conn.(*net.TCPConn); ok {
		_ = c.SetKeepAlive(true)
	}

	return conn, err
}

func (session *Session) closeConn(ctx context.Context, conn net.Conn) {
	if conn == nil {
		return
	}
	if session.verbos {
		log.CtxWarn(ctx, "closeConn %s", conn.LocalAddr().String())
	}
	if err := conn.Close(); err != nil {
		log.CtxError(ctx, "closeConn fail %s:%s", conn.LocalAddr().String(), err.Error())
	}
}
