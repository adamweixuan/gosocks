package main

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"time"
)

var (
	authReply = []byte{
		socksVersion,
		noAuthRequired,
	}

	authFailReply = []byte{
		socksVersion,
		noAcceptableAuth,
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
		_, _ = session.local.Write(authFailReply)
		session.closeConn(ctx, session.local)
		return
	}

	if session.verbos {
		log.CtxTrace(ctx, "auth %s success", session.local.LocalAddr().String())
	}

	msg := session.connect(ctx)

	_, _ = session.local.Write(NewReply(msg.Code()))

	if msg.Code() != success {
		log.CtxError(ctx, "connect fail:%s", msg.Err())
		return
	}

	if session.verbos {
		log.CtxTrace(ctx, "connectCmd to %s success", session.remote.RemoteAddr().String())
	}
	session.relay(ctx)
}

func (session *Session) auth() error {
	var buf [2]byte

	if _, err := io.ReadFull(session.local, buf[:]); err != nil {
		return err
	}
	ver, count := int(buf[0]), int(buf[1])
	if ver != socksVersion {
		return ErrUnSupportVersion
	}

	methods := make([]byte, count)

	if _, err := io.ReadFull(session.local, methods); err != nil {
		return err
	}

	// todo auth
	_, err := session.local.Write(authReply)
	return err
}

func (session *Session) connect(_ context.Context) *Message {
	var buf [4]byte
	if _, err := io.ReadFull(session.local, buf[:]); err != nil {
		return NewMessage(readConnErr, err)
	}

	cmd, atyp := cmdType(buf[1]), addrtype(buf[3])

	if cmd == bindCmd || cmd == udpAssociateCmd {
		return NewMessage(cmdNotSupported, ErrUnsupportCmd)
	}

	var addr string

	switch atyp {
	case ipv4:
		var ip [4]byte
		if _, err := io.ReadFull(session.local, ip[:]); err != nil {
			return NewMessage(readConnErr, err)
		}
		addr = net.IP(ip[:]).String()
	case domain:
		var cnt [1]byte
		if _, err := io.ReadFull(session.local, cnt[:]); err != nil {
			return NewMessage(readConnErr, err)
		}

		domainSize := int(cnt[0])
		buf := make([]byte, domainSize)

		if _, err := io.ReadFull(session.local, buf[:]); err != nil {
			return NewMessage(readConnErr, err)
		}
		addr = string(buf)
	case ipv6:
		var ip [16]byte
		if _, err := io.ReadFull(session.local, ip[:]); err != nil {
			return NewMessage(readConnErr, err)
		}
		addr = net.IP(ip[:]).String()
	default:
		return NewMessage(addrNotSupported, ErrUnsupportAddrType)
	}

	var portBuf [2]byte

	if _, err := io.ReadFull(session.local, portBuf[:]); err != nil {
		return NewMessage(readConnErr, err)
	}

	port := binary.BigEndian.Uint16(portBuf[:])

	endPoint := net.JoinHostPort(addr, strconv.Itoa(int(port)))
	dst, err := session.dial(endPoint)
	if err != nil {
		return NewMessage(readConnErr, err)
	}
	session.remote = dst
	return NewMessage(success, nil)
}

func (session *Session) relay(ctx context.Context) {
	forward := func(from, to net.Conn) {
		buf := Get(defaultBufSize, defaultBufSize)
		defer Put(buf)
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
