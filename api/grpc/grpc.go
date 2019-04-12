package grpc

import (
	"net"

	"coretest-go/api"
	"google.golang.org/grpc"
)

type Server struct {
	server  *grpc.Server
	user    api.UserService
	order   api.OrderService
	payment api.PaymentService
	logistic api.LogisticService
}

// Start grpc service
func (s *Server) Serve(lis net.Listener) error {
	s.server = grpc.NewServer()
	return s.server.Serve(lis)
}

func (s *Server) Shutdown() error {
	s.server.GracefulStop()
	return nil
}
