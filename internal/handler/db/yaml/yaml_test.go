package yaml

import (
	"testing"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGetSeeds(t *testing.T) {
	t.Parallel()

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
			Labels: map[string]string{
				"label-1": "value-1",
				"label-2": "value-2",
			},
		},
	}

	got, err := GetSeeds("./testdata/seed.yaml")
	if err != nil {
		t.Errorf("failed to get seeds: %v", err)
		return
	}

	if diff := cmp.Diff(got, want, protocmp.IgnoreFields(&runpb.Service{}, "create_time"), protocmp.Transform()); diff != "" {
		t.Errorf("\n(-got, +want)\n%s", diff)
		return
	}
}
