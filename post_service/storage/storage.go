package storage

import (
	"4microservice/post-service/storage/postgres"
	"4microservice/post-service/storage/repo"
	"database/sql"
)

type IStorage interface {
	Post() repo.PostStorageI
}

type Pg struct {
	db       *sql.DB
	postRepo repo.PostStorageI
}

func NewStoragePg(db *sql.DB) *Pg {
	return &Pg{
		db:       db,
		postRepo: postgres.NewPostRepo(db),
	}
}

func (s Pg) Post() repo.PostStorageI {
	return s.postRepo
}
