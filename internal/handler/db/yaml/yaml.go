package yaml

import (
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/goccy/go-yaml"
)

type seeds struct {
	Services []*runpb.Service `yaml:"services"`
}

func GetSeeds(path string) (_ []*runpb.Service, reterr error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the seed file: %w", err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			reterr = fmt.Errorf("failed to close the seed file: %w", err)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read the seed file: %w", err)
	}

	var s seeds
	if err := yaml.Unmarshal(content, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the seed file: %w", err)
	}

	return s.Services, nil
}
