package memory

import (
	"context"
	"github.com/SoroushBeigi/movie-microservice/rating/internal/repository"
	"github.com/SoroushBeigi/movie-microservice/rating/pkg/model"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

func (r *Repository) Get(ctx context.Context, recordId model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if ratings, ok := r.data[recordType][recordId]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordId], nil
}

func (r *Repository) Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.
		Rating{}
	}
	r.data[recordType][recordID] =
		append(r.data[recordType][recordID], *rating)
	return nil
}
