package mysql

import "database/sql"

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
