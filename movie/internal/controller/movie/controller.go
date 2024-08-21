package movie

import (
	"context"
	"errors"
	metadatamodel "github.com/SoroushBeigi/movie-microservice/metadata/pkg/model"
	"github.com/SoroushBeigi/movie-microservice/movie/internal/gateway"
	"github.com/SoroushBeigi/movie-microservice/movie/pkg/model"
	ratingmodel "github.com/SoroushBeigi/movie-microservice/rating/pkg/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) (float64, error)
	PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

type metadataGateway interface {
	Get(ctx context.Context, id string) (*metadatamodel.
		Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func New(ratingGateway ratingGateway, metadataGateway metadataGateway) *Controller {
	return &Controller{ratingGateway, metadataGateway}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	details := &model.MovieDetails{Metadata: *metadata}

	rating, ratingErr := c.ratingGateway.GetAggregatedRating(ctx,
		ratingmodel.RecordID(id), ratingmodel.RecordTypeMovie)

	if ratingErr != nil && !errors.Is(ratingErr, gateway.ErrNotFound) {
	} else if ratingErr != nil {
		return nil, ratingErr
	} else {
		details.Rating = &rating
	}
	return details, nil
}
