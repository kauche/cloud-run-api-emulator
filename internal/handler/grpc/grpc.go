package grpc

import (
	"context"
	"fmt"
	"net"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/110y/servergroup"
	"google.golang.org/grpc"

	"github.com/kauche/cloud-run-api-emulator/internal/usecase"
)

var (
	_ servergroup.Server  = (*Server)(nil)
	_ servergroup.Stopper = (*Server)(nil)
)

type Server struct {
	server *grpc.Server
	port   int
}

func NewServer(port int, servicesUsecase *usecase.ServicesUsecase) *Server {
	server := grpc.NewServer()

	runpb.RegisterServicesServer(server, &servicesServer{uc: servicesUsecase})
	runpb.RegisterRevisionsServer(server, &revisionsServer{})
	runpb.RegisterJobsServer(server, &jobsServer{})
	runpb.RegisterTasksServer(server, &tasksServer{})
	runpb.RegisterExecutionsServer(server, &executionsServer{})

	return &Server{
		server: server,
		port:   port,
	}
}

func (s *Server) Start(_ context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on the port %d: %w", s.port, err)
	}

	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("the server aborted: %w", err)
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()

	return nil
}
