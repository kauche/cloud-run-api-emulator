package domain

import (
	"context"

	"cloud.google.com/go/run/apiv2/runpb"
)

type ServicesRepository interface {
	CreateService(ctx context.Context, parent string, service *runpb.Service) error
	CreateServices(ctx context.Context, service []*runpb.Service) error
	ListServices(ctx context.Context, parent string, limit int32) ([]*runpb.Service, error)
	ListServicesByParentCreatedAtName(ctx context.Context, parent, createdAt, name string, limit int32) ([]*runpb.Service, error)
}
