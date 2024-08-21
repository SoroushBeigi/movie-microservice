package model

import (
	"github.com/SoroushBeigi/movie-microservice/metadata/pkg/model"
)

type MovieDetails struct {
	Rating   *float64       `json:"rating,omitempty"`
	Metadata model.Metadata `json:"metadata"`
}
