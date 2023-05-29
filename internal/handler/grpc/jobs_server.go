package grpc

import "cloud.google.com/go/run/apiv2/runpb"

type jobsServer struct {
	runpb.UnimplementedJobsServer
}
