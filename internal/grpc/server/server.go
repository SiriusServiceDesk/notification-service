package server

import (
	"net"

	_ "github.com/jnewmano/grpc-json-proxy/codec"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	listener   net.Listener
	grpcServer *grpc.Server
}

func NewGRPCServer(listener net.Listener, server *grpc.Server) *GRPCServer {
	return &GRPCServer{
		listener:   listener,
		grpcServer: server,
	}
}

func (s *GRPCServer) Start() error {
	return s.grpcServer.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() {
	s.grpcServer.GracefulStop()
}
