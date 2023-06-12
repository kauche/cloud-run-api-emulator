package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/110y/run"
	"github.com/110y/servergroup"

	"github.com/kauche/cloud-run-api-emulator/internal/handler/db/sqlite"
	"github.com/kauche/cloud-run-api-emulator/internal/handler/db/yaml"
	"github.com/kauche/cloud-run-api-emulator/internal/handler/grpc"
	"github.com/kauche/cloud-run-api-emulator/internal/usecase"
)

func Run() {
	run.Run(server)
}

func server(ctx context.Context) (code int) {
	data := flag.String("data", "", "A database file path to persist the data")
	seed := flag.String("seed", "", "A path to a YAML file that contains the seed data")

	flag.Parse()

	if *data == "" {
		*data = ":memory:"
	}

	db, err := sqlite.NewDB(*data)
	if err != nil {
		// TODO: structured log
		fmt.Fprintf(os.Stderr, "failed to create the database connection: %s\n", err)
		return 1
	}
	defer func() {
		if err := db.Close(); err != nil {
			// TODO: structured log
			fmt.Fprintf(os.Stderr, "failed to close the database connection: %s\n", err)
			code = 1
		}
	}()

	srepo := sqlite.NewServicesRepository(db)
	suc := usecase.NewServicesUsecase(srepo)

	gs := grpc.NewServer(8000, suc) // TODO: make the port configurable

	if *seed != "" {
		seedServices, err := yaml.GetSeeds(*seed)
		if err != nil {
			// TODO: structured log
			fmt.Fprintf(os.Stderr, "failed to read the seed file: %s\n", err)
			return 0
		}

		if err := srepo.CreateServices(ctx, seedServices); err != nil {
			// TODO: structured log
			fmt.Fprintf(os.Stderr, "failed to save the seed services: %s\n", err)
			return 0
		}
	}

	var sg servergroup.Group
	sg.Add(gs)

	if err := sg.Start(ctx); err != nil {
		// TODO: structured log
		fmt.Fprintf(os.Stderr, "the server aborted: %s\n", err)
		return 1
	}

	return 0
}
