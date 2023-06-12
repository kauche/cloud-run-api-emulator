package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

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
		UID:         service.Uid,
		Generation:  service.Generation,
		CreatedAt:   xo.NewTime(service.CreateTime.AsTime()),
	}

	if err := s.Insert(ctx, r.db); err != nil {
		return fmt.Errorf("failed to save the service: %w", err)
	}

	// TODO: should insert by bulk
	for k, v := range service.Annotations {
		sa := &xo.ServiceAnnotation{
			ServiceParent: parent,
			ServiceName:   service.Name,
			Key:           k,
			Value:         v,
		}

		if err := sa.Insert(ctx, r.db); err != nil {
			return fmt.Errorf("failed to save the service annotation: %w", err)
		}
	}

	// TODO: should insert by bulk
	for k, v := range service.Labels {
		sl := &xo.ServiceLabel{
			ServiceParent: parent,
			ServiceName:   service.Name,
			Key:           k,
			Value:         v,
		}

		if err := sl.Insert(ctx, r.db); err != nil {
			return fmt.Errorf("failed to save the service label: %w", err)
		}
	}

	return nil
}

func (r *ServicesRepository) CreateServices(ctx context.Context, services []*runpb.Service) error {
	// TODO: should insert by bulk
	for _, s := range services {
		nameParts := strings.Split(s.Name, "/")
		if len(nameParts) != 6 {
			return fmt.Errorf("invalid service name: %s", s.Name)
		}
		parent := fmt.Sprintf("projects/%s/locations/%s", nameParts[1], nameParts[3])

		if err := r.CreateService(ctx, parent, s); err != nil {
			return fmt.Errorf("failed to insert service, name=%s : %w", s.Name, err)
		}
	}

	return nil
}

// TODO: annotation
func (r *ServicesRepository) ListServices(ctx context.Context, parent string, limit int32) ([]*runpb.Service, error) {
	res, err := xo.ListServicesByParentLimit(ctx, r.db, parent, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}

	services := make([]*runpb.Service, len(res))
	for i, s := range res {
		services[i] = &runpb.Service{
			Name:        s.Name,
			Description: s.Description,
			Uid:         s.UID,
			Uri:         s.URI,
			Generation:  s.Generation,
			CreateTime:  timestamppb.New(s.CreatedAt.Time()),
		}

		services[i].Annotations = make(map[string]string)
		annotations, err := r.listServiceAnnotations(ctx, s.Parent, s.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to list service annnotations: %w", err)
		}

		for _, a := range annotations {
			services[i].Annotations[a.Key] = a.Value
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
			Name:        s.Name,
			Description: s.Description,
			Uid:         s.UID,
			Uri:         s.URI,
			Generation:  s.Generation,
			CreateTime:  timestamppb.New(s.CreatedAt.Time()),
		}

		services[i].Annotations = make(map[string]string)
		annotations, err := r.listServiceAnnotations(ctx, s.Parent, s.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to list service annnotations: %w", err)
		}

		for _, a := range annotations {
			services[i].Annotations[a.Key] = a.Value
		}
	}

	return services, nil
}

func (s *ServicesRepository) listServiceAnnotations(ctx context.Context, parent, name string) ([]*xo.ServiceAnnotation, error) {
	res, err := xo.ListServiceAnnotationsByParentName(ctx, s.db, parent, name)
	if err != nil {
		return nil, fmt.Errorf("failed to list service annotations: %w", err)
	}

	annotations := make([]*xo.ServiceAnnotation, len(res))
	for i, a := range res {
		annotations[i] = &xo.ServiceAnnotation{
			ServiceParent: a.ServiceParent,
			ServiceName:   a.ServiceName,
			Key:           a.Key,
			Value:         a.Value,
		}
	}

	return annotations, nil
}
