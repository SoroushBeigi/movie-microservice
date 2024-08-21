package model

import model "github.com/SoroushBeigi/movie-microservice/metadeta/pkg"

type MovieDetails struct {
	Rating   *float64       `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
