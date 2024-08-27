package main

import (
	"context"
	"fmt"
	"github.com/SoroushBeigi/movie-microservice/gen"
	"github.com/SoroushBeigi/movie-microservice/movie/internal/controller/movie"
	metadatagateway "github.com/SoroushBeigi/movie-microservice/movie/internal/gateway/metadata/http"
	ratinggateway "github.com/SoroushBeigi/movie-microservice/movie/internal/gateway/rating/http"
	grpchandler "github.com/SoroushBeigi/movie-microservice/movie/internal/handler/grpc"
	"github.com/SoroushBeigi/movie-microservice/pkg/discovery"
	"github.com/SoroushBeigi/movie-microservice/pkg/discovery/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"time"
)

const serviceName = "movie"

func main() {
	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	var cfg config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	port := cfg.API.Port
	log.Printf("Starting the metadata service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.
		GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID,
		serviceName, fmt.Sprintf("localhost:%d", port)); err !=
		nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)

	ctrl := movie.NewController(ratingGateway, metadataGateway)
	h := grpchandler.New(ctrl)

	lis, listenErr := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if listenErr != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMovieServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
