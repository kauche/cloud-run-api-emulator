package grpc

import "cloud.google.com/go/run/apiv2/runpb"

type revisionsServer struct {
	runpb.UnimplementedRevisionsServer
}
