package grpc

import (
	"context"
	"fmt"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kauche/cloud-run-api-emulator/internal/usecase"
)

type servicesServer struct {
	uc *usecase.ServicesUsecase
	runpb.UnimplementedServicesServer
}

func (s *servicesServer) CreateService(ctx context.Context, req *runpb.CreateServiceRequest) (*longrunningpb.Operation, error) {
	if err := s.uc.CreateService(ctx, req); err != nil {
		// TODO: check if the request is valid
		// TODO: check if already exists

		// TODO: log error
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create the service: %s", err))
	}

	return &longrunningpb.Operation{
		// TODO: fill other fields
		Done: true,
	}, nil
}

func (s *servicesServer) ListServices(ctx context.Context, req *runpb.ListServicesRequest) (*runpb.ListServicesResponse, error) {
	services, nextPageToken, err := s.uc.ListServices(ctx, req)
	if err != nil {
		// TODO: log error
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to list services: %s", err))
	}

	return &runpb.ListServicesResponse{
		Services:      services,
		NextPageToken: nextPageToken,
	}, nil
}
