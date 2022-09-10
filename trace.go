package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fastrand"
)

const (
	length     = 30
	KeyTraceID = "TRACEID"
	min        = 1 << 20
	max        = 1 << 30
)

func ms() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func random(min, max uint32) uint32 {
	return fastrand.Uint32n(max+min) - min
}

func genTraceID() string {
	sb := strings.Builder{}
	sb.Grow(length)
	sb.WriteString(strconv.FormatUint(uint64(ms()), 16))
	sb.WriteString(strconv.FormatUint(uint64(ip2int(localIP)), 16))
	sb.WriteString(strconv.FormatUint(uint64(random(min, max)), 16))

	return sb.String()
}

func NewCtxWithTraceID() context.Context {
	return context.WithValue(context.Background(), "KeyTraceID", genTraceID()) //nolint:staticcheck
}
