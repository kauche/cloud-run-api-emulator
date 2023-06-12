# Cloud Run API Emulator

#### :warning: This project is still in the super experimental phase :warning:

Cloud Run API Emulator provides application developers with a locally-running, emulated instance of Cloud Run API to enable local development and testing.

## Usage

Run cloud-run-api-emulator by using docker like below:

```
$ docker run --publish 8000:8000 --detach ghcr.io/kauche/cloud-run-api-emulator:0.0.3
```

And then, you can use the emulator through a gRPC client like below:

(The example below is written in Go, but you can use any language.)

```go
package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx := context.Background()

	cc, err := grpc.DialContext(ctx, "localhost:8000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to the emulator: %s\n", err)
		os.Exit(1)
	}
	defer func() {
		if err = cc.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close the connection: %s\n", err)
			os.Exit(1)
		}
	}()

	client := runpb.NewServicesClient(cc)

	createReq := &runpb.CreateServiceRequest{
		Parent:    "projects/test-project/locations/us-central1",
		ServiceId: "my-service",
		Service: &runpb.Service{
			Description: "my service",
		},
	}

	_, err = client.CreateService(ctx, createReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create the service: %s\n", err)
		os.Exit(1)
	}

	listReq := &runpb.ListServicesRequest{
		Parent:   "projects/test-project/locations/us-central1",
		PageSize: 5,
	}

	listRes, err := client.ListServices(ctx, listReq)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list services: %s\n", err)
		os.Exit(1)
	}

	// this prints `service: projects/test-project/locations/us-central1/services/my-service`
	for _, service := range listRes.Services {
		fmt.Printf("service: %s\n", service.Name)
	}
}

```

## Supported Methods

### Services

-   [x] CreateService
    -   ValidateOnly mode is not yet supported.
    -   Some fields of Service are not yet supported.
-   [ ] GetService
-   [x] ListServices
    -   ShowDeleted is not yet supported.
    -   Some fields of Service are not yet supported.
-   [ ] UpdateService
-   [ ] DeleteService
-   [ ] GetIamPolicy
-   [ ] SetIamPolicy
-   [ ] TestIamPermissions

### Revisions

-   [ ] GetRevision
-   [ ] ListRevisions
-   [ ] DeleteRevision

### Jobs

-   [ ] CreateJob
-   [ ] GetJob
-   [ ] ListJobs
-   [ ] UpdateJob
-   [ ] DeleteJob
-   [ ] RunJob
-   [ ] GetIamPolicy
-   [ ] SetIamPolicy
-   [ ] TestIamPermissions

### Tasks

-   [ ] GetTask
-   [ ] ListTasks

### Executions

-   [ ] GetExecution
-   [ ] ListExecutions
-   [ ] DeleteExecution
