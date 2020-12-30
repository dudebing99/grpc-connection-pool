# Usage
```bash
go get github.com/dudebing99/grpc-connection-pool
```

# Example
```go
package main

import (
	"context"
	"fmt"
	rpc "github.com/dudebing99/grpc-connection-pool"
)

func main() {
	pool, err := rpc.NewRpcClientPool(rpc.WithServerAddr("0.0.0.0:8080"))
	if err != nil {
		fmt.Println("init client pool error")
		return
	}

	clientConn, close, err := pool.Acquire()
	defer close()
	if err != nil {
		fmt.Println("acquire client connection error")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	reply, err := NewGreeterClient(clientConn).SayHello(ctx, &HelloRequest{Name: "SillyBoy"})
	if err != nil {
		fmt.Println("say hello error, ", err)
		return
	}

	fmt.Println(reply.Message)
}
```
