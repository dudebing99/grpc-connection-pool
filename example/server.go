package main

import (
	"fmt"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	return &HelloReply{
		Message: "hello, " + in.Name,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	server := grpc.NewServer()
	RegisterGreeterServer(server, &Server{})

	server.Serve(lis)
}
