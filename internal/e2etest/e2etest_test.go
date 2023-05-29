package e2etest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"cloud.google.com/go/run/apiv2/runpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx := context.Background()

		target := os.Getenv("EMULATOR_TARGET")
		if target == "" {
			fmt.Fprintln(os.Stderr, "the emulator target is not set")
			return 1
		}

		cc, err := grpc.DialContext(ctx, target, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to connect to the emulator: %s\n", err)
			return 1
		}

		serviceClient = runpb.NewServicesClient(cc)

		return m.Run()
	}())
}
