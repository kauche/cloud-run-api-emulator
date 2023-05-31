package e2etest

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"google.golang.org/protobuf/testing/protocmp"

	"github.com/kauche/cloud-run-api-emulator/internal/handler/db/sqlite"
)

var serviceClient runpb.ServicesClient

func TestServices(t *testing.T) {
	t.Parallel()

	project := uuid.New().String()
	parent := fmt.Sprintf("projects/%s/locations/us-central1", project)

	t.Cleanup(func() {
		db, err := sqlite.NewDB("../../bin/cloud-run-api-emulator.db")
		if err != nil {
			t.Errorf("failed to create the database connection for the clean up: %s", err)
			return
		}
		defer db.Close()

		_, err = db.Exec("DELETE FROM services WHERE parent = $1", parent)
		if err != nil {
			t.Errorf("failed to delete services for the celan up: %s", err)
			return
		}
	})

	ctx := context.Background()

	numServices := 10
	pageSize := 2

	for i := 1; i <= numServices; i++ {
		req := &runpb.CreateServiceRequest{
			Parent:    parent,
			ServiceId: fmt.Sprintf("test-service-%d", i),
			Service: &runpb.Service{
				Description: fmt.Sprintf("test service %d", i),
				Uri:         fmt.Sprintf("service-%d.example.com", i),
				Annotations: map[string]string{
					"annotation-key-1": "annotation-value-1",
					"annotation-key-2": "annotation-value-2",
				},
			},
		}

		got, err := serviceClient.CreateService(ctx, req)
		if err != nil {
			t.Errorf("failed to create service: %s", err)
			return
		}

		want := &longrunningpb.Operation{
			Done: true,
		}

		if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
			t.Errorf("\n(-got, +want)\n%s", diff)
			return
		}
	}

	var pageToken string
	for i := 0; i < numServices/pageSize; i++ {
		req := &runpb.ListServicesRequest{
			Parent:    parent,
			PageSize:  int32(pageSize),
			PageToken: pageToken,
		}

		res, err := serviceClient.ListServices(ctx, req)
		if err != nil {
			t.Errorf("failed to list services: %s", err)
			return
		}

		serviceNumber := (numServices/pageSize - i) * pageSize
		want := []*runpb.Service{
			{
				Name:        fmt.Sprintf("%s/services/test-service-%d", parent, serviceNumber),
				Description: fmt.Sprintf("test service %d", serviceNumber),
				Uri:         fmt.Sprintf("service-%d.example.com", serviceNumber),
				Annotations: map[string]string{
					"annotation-key-1": "annotation-value-1",
					"annotation-key-2": "annotation-value-2",
				},
			},
			{
				Name:        fmt.Sprintf("%s/services/test-service-%d", parent, serviceNumber-1),
				Description: fmt.Sprintf("test service %d", serviceNumber-1),
				Uri:         fmt.Sprintf("service-%d.example.com", serviceNumber-1),
				Annotations: map[string]string{
					"annotation-key-1": "annotation-value-1",
					"annotation-key-2": "annotation-value-2",
				},
			},
		}

		if diff := cmp.Diff(res.Services, want, protocmp.IgnoreFields(&runpb.Service{}, "create_time"), protocmp.Transform()); diff != "" {
			t.Errorf("\n(-got, +want)\n%s", diff)
			return
		}

		pageToken = res.GetNextPageToken()
	}

	if pageToken != "" {
		t.Errorf("pageToken should be empty at the last page, but got %s", pageToken)
	}
}

func TestServices_Seed(t *testing.T) {
	ctx := context.Background()

	req := &runpb.ListServicesRequest{
		Parent:   "projects/test-project/locations/us-central1",
		PageSize: 5,
	}

	res, err := serviceClient.ListServices(ctx, req)
	if err != nil {
		t.Errorf("failed to list services: %s", err)
		return
	}

	want := []*runpb.Service{
		{
			Name: "projects/test-project/locations/us-central1/services/service-1",
			Annotations: map[string]string{
				"annotation-1": "value-1",
				"annotation-2": "value-2",
			},
		},
		{
			Name: "projects/test-project/locations/us-central1/services/service-2",
			Annotations: map[string]string{
				"annotation-1": "value-1",
				"annotation-2": "value-2",
			},
		},
	}

	if diff := cmp.Diff(res.Services, want, protocmp.IgnoreFields(&runpb.Service{}, "create_time"), protocmp.Transform()); diff != "" {
		t.Errorf("\n(-got, +want)\n%s", diff)
		return
	}
}
