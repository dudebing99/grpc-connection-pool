package rpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync/atomic"
	"time"
)

func connect(addr string, dynamicLink bool, ops ...grpc.DialOption) (*grpc.ClientConn, error) {
	if dynamicLink == true {
		return grpc.Dial(addr, ops...)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return grpc.DialContext(ctx, addr, ops...)
}

func (pool *Pool) init() error {
	len := cap(pool.connections) - len(pool.connections)
	addr := pool.ServerAddr
	dynamicLink := pool.DynamicLink
	ops := pool.dialOptions

	for i := 1; i <= len; i++ {
		client, err := connect(addr, dynamicLink, ops...)
		if err != nil {
			return err
		}
		pool.connections <- client
	}

	return nil
}

func (pool *Pool) count(add int64) {
	atomic.AddInt64(&pool.counter, add)
}

// close Recycling available links
func (pool *Pool) close(conn *grpc.ClientConn) {
	// double check
	if conn == nil {
		return
	}

	go func() {
		detect, _ := passivate(conn)
		if detect && pool.ChannelStat {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			select {
			case pool.connections <- conn:
			case <-ctx.Done():
				destroy(conn)
			}
		}
		pool.count(-1)
	}()
}

// destroy tears down the ClientConn and all underlying connections.
func destroy(conn *grpc.ClientConn) error {
	return conn.Close()
}

type condition = int

const (
	// Ready Can be used
	Ready condition = iota
	// Put Not available. Maybe later.
	Put
	// Destroy Failure occurs and cannot be restored
	Destroy
)

// passivate Action before releasing the resource
func passivate(conn *grpc.ClientConn) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if conn.WaitForStateChange(ctx, connectivity.Ready) && conn.WaitForStateChange(ctx, connectivity.Shutdown) && conn.WaitForStateChange(ctx, connectivity.Idle) {
		return true, nil
	}

	return false, destroy(conn)
}

// activate Action taken after getting the resource
func activate(conn *grpc.ClientConn) int {
	stat := conn.GetState()
	switch {
	case stat == connectivity.Ready:
		return Ready
	case stat == connectivity.Shutdown:
		return Destroy
	default:
		return Put
	}
}
