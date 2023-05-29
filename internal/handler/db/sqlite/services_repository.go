package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kauche/cloud-run-api-emulator/internal/domain"
	"github.com/kauche/cloud-run-api-emulator/internal/handler/db/sqlite/xo"
)

var _ domain.ServicesRepository = (*ServicesRepository)(nil)

func NewServicesRepository(db *sql.DB) *ServicesRepository {
	return &ServicesRepository{
		db: db,
	}
}

type ServicesRepository struct {
	db *sql.DB
}

func (r *ServicesRepository) CreateService(ctx context.Context, parent string, service *runpb.Service) error {
	s := &xo.Service{
		Parent:      parent,
		Name:        service.Name,
		Description: service.Description,
		URI:         service.Uri,
		CreatedAt:   xo.NewTime(service.CreateTime.AsTime()),
	}

	if err := s.Insert(ctx, r.db); err != nil {
		return fmt.Errorf("failed to save the service: %w", err)
	}

	// TODO: should insert by bulk
	for k, v := range service.Annotations {
		sa := &xo.ServiceAnnotation{
			ServiceName: service.Name,
			Key:         k,
			Value:       v,
		}

		if err := sa.Insert(ctx, r.db); err != nil {
			return fmt.Errorf("failed to save the service annotation: %w", err)
		}
	}

	// TODO: should insert by bulk
	for k, v := range service.Labels {
		sl := &xo.ServiceLabel{
			ServiceName: service.Name,
			Key:         k,
			Value:       v,
		}

		if err := sl.Insert(ctx, r.db); err != nil {
			return fmt.Errorf("failed to save the service label: %w", err)
		}
	}

	return nil
}

func (r *ServicesRepository) ListServices(ctx context.Context, parent string, limit int32) ([]*runpb.Service, error) {
	res, err := xo.ListServicesByParentLimit(ctx, r.db, parent, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := make([]*runpb.Service, len(res))
	for i, s := range res {
		services[i] = &runpb.Service{
			Name:        fmt.Sprintf("%s/services/%s", s.Parent, s.Name),
			Description: s.Description,
			Uid:         s.UID,
			Generation:  int64(s.Generation),
			CreateTime:  timestamppb.New(s.CreatedAt.Time()),
		}
	}

	return services, nil
}

func (r *ServicesRepository) ListServicesByParentCreatedAtName(ctx context.Context, parent string, createdAt, name string, limit int32) ([]*runpb.Service, error) {
	res, err := xo.ListNextServicesByParentCreatedAtNameLimit(ctx, r.db, parent, createdAt, name, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := make([]*runpb.Service, len(res))
	for i, s := range res {
		services[i] = &runpb.Service{
			Name:        fmt.Sprintf("%s/services/%s", s.Parent, s.Name),
			Description: s.Description,
			Uid:         s.UID,
			Generation:  int64(s.Generation),
			CreateTime:  timestamppb.New(s.CreatedAt.Time()),
		}
	}

	return services, nil
}
