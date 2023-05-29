package grpc

import "cloud.google.com/go/run/apiv2/runpb"

type executionsServer struct {
	runpb.UnimplementedExecutionsServer
}
