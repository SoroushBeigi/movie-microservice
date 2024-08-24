package grpc

import (
	"context"
	"github.com/SoroushBeigi/movie-microservice/gen"
	"github.com/SoroushBeigi/movie-microservice/internal/grpcutil"
	"github.com/SoroushBeigi/movie-microservice/metadata/pkg/model"
	"github.com/SoroushBeigi/movie-microservice/pkg/discovery"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, connErr := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if connErr != nil {
		return nil, connErr
	}
	defer conn.Close()
	client := gen.NewMetadataServiceClient(conn)
	resp, getErr := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if getErr != nil {
		return nil, getErr
	}
	return model.MetadataFromProto(resp.Metadata), nil
}
