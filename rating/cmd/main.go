package main

import (
	"context"
	"fmt"
	"github.com/SoroushBeigi/movie-microservice/gen"
	"github.com/SoroushBeigi/movie-microservice/pkg/discovery"
	"github.com/SoroushBeigi/movie-microservice/pkg/discovery/consul"
	"github.com/SoroushBeigi/movie-microservice/rating/internal/controller/rating"
	grpchandler "github.com/SoroushBeigi/movie-microservice/rating/internal/handler/grpc"
	"github.com/SoroushBeigi/movie-microservice/rating/internal/repository/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v3"
	"log"
	"net"
	"os"
	"time"
)

const serviceName = "rating"

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

	registry, consulErr := consul.NewRegistry("localhost:8500")
	if consulErr != nil {
		panic(consulErr)
	}
	ctx := context.Background()
	instanceID := discovery.
		GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
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

	repo, sqlErr := mysql.New()
	if sqlErr != nil {
		panic(sqlErr)
	}
	ctrl := rating.NewController(repo, nil)
	h := grpchandler.New(ctrl)
	lis, lisErr := net.Listen("tcp", fmt.Sprintf("localhost:%v", port))
	if lisErr != nil {
		log.Fatalf("failed to listen: %v", lisErr)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)
	if err := srv.Serve(lis); err != nil {
		panic(err)
	}
}
