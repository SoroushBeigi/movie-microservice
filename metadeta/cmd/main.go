package cmd

import (
	"github.com/SoroushBeigi/movie-microservice/metadeta/internal/controller/metadata"
	httphandler "github.com/SoroushBeigi/movie-microservice/metadeta/internal/handler/http"
	"github.com/SoroushBeigi/movie-microservice/metadeta/internal/repository/memory"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
