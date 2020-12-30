package rpc

import (
	"google.golang.org/grpc"
	"sync"
	"time"
)

type chanStat = bool

type Pool struct {
	ServerAddr     string
	MaxCap         int64
	AcquireTimeout time.Duration
	DynamicLink    bool
	OverflowCap    bool
	dialOptions    []grpc.DialOption

	lock        *sync.Mutex
	connections chan *grpc.ClientConn
	ChannelStat chanStat
	counter     int64
}

type Option func(*Pool)

func WithMaxCap(num int64) Option {
	return func(pool *Pool) {
		pool.MaxCap = num
	}
}

func WithServerAddr(addr string) Option {
	return func(pool *Pool) {
		pool.ServerAddr = addr
	}
}

func WithAcquireTimeOut(timeout time.Duration) Option {
	return func(pool *Pool) {
		pool.AcquireTimeout = timeout
	}
}

func WithOverflowCap(overflowCap bool) Option {
	return func(pool *Pool) {
		pool.OverflowCap = overflowCap
	}
}

func WithDynamicLink(dynamicLink bool) Option {
	return func(pool *Pool) {
		pool.DynamicLink = dynamicLink
	}
}

func WithDialOption(ops ...grpc.DialOption) Option {
	return func(pool *Pool) {
		pool.dialOptions = ops
	}
}
