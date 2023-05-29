package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/kauche/bjt"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kauche/cloud-run-api-emulator/internal/domain"
)

var ErrEmptyService = errors.New("service is nil")

type listServicesPageToken struct {
	CreatedAt string
	Name      string
}

func NewServicesUsecase(repo domain.ServicesRepository) *ServicesUsecase {
	return &ServicesUsecase{
		repo: repo,
	}
}

type ServicesUsecase struct {
	repo domain.ServicesRepository
}

func (u *ServicesUsecase) CreateService(ctx context.Context, req *runpb.CreateServiceRequest) error {
	service := req.GetService()
	if service == nil {
		return ErrEmptyService
	}

	service.Name = req.ServiceId
	service.CreateTime = timestamppb.Now()

	if err := u.repo.CreateService(ctx, req.Parent, service); err != nil {
		// TODO: check if already exists
		return fmt.Errorf("failed to persist the service: %w", err)
	}

	return nil
}

func (u *ServicesUsecase) ListServices(ctx context.Context, req *runpb.ListServicesRequest) ([]*runpb.Service, string, error) {
	limit := req.GetPageSize()
	if limit == 0 {
		limit = 100 // TODO: set the default value that is used in the actual Cloud Run API
	}

	// NOTE: Fetch one more service to check if there are more services for the next page
	limit++

	var services []*runpb.Service

	if req.PageToken == "" {
		res, err := u.repo.ListServices(ctx, req.Parent, limit)
		if err != nil {
			return nil, "", fmt.Errorf("failed to list services: %w", err)
		}

		services = res
	} else {
		token, err := bjt.Decode[listServicesPageToken](req.PageToken)
		if err != nil {
			return nil, "", fmt.Errorf("failed to decode the page token: %w", err)
		}

		res, err := u.repo.ListServicesByParentCreatedAtName(ctx, req.Parent, token.Source.CreatedAt, token.Source.Name, limit)
		if err != nil {
			return nil, "", fmt.Errorf("failed to list services by parent, created_at and name: %w", err)
		}

		services = res
	}

	numFetchedServices := len(services)

	if numFetchedServices == 0 {
		return []*runpb.Service{}, "", nil
	} else if int32(numFetchedServices) < limit {
		return services, "", nil
	}

	lastService := services[numFetchedServices-2]

	// NOTE: Service.Name consists of projects/{project}/locations/{location}/services/{service_id}
	nameParts := strings.Split(lastService.Name, "/")
	if len(nameParts) != 6 {
		return nil, "", fmt.Errorf("invalid service name: %s", lastService.Name)
	}
	lastServiceName := nameParts[5]

	pageToken := bjt.NewToken(&listServicesPageToken{
		CreatedAt: lastService.CreateTime.AsTime().Format("2006-01-02 15:04:05.999999999"),
		Name:      lastServiceName,
	})

	token, err := pageToken.Encode()
	if err != nil {
		return nil, "", fmt.Errorf("failed to encode the page token: %w", err)
	}

	return services[:numFetchedServices-1], token, nil
}
