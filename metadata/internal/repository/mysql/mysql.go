package mysql

import (
	"context"
	"database/sql"
	"github.com/SoroushBeigi/movie-microservice/metadata/internal/repository"
	"github.com/SoroushBeigi/movie-microservice/metadata/pkg/model"
)

type Repository struct {
	db *sql.DB
}

func New() (repository *Repository, err error) {
	//storing database credentials in code is generally considered a bad practice and is done for testing purposes here!
	db, err := sql.Open("mysql", "root:password@/movieexample")
	if err != nil {
		return nil, err
	}
	return &Repository{db}, nil
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "SELECT title,description,director FROM movies WHERE id = ?", id)
	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &model.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

func (r *Repository) Put(ctx context.Context, id string, metadata *model.Metadata) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO movies (id,title,description,director) VALUES (?,?,?,?)",
		id, metadata.Title, metadata.Description, metadata.Director)
	return err
}
