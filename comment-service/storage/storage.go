package storage

import (
	"4microservices/comment-service/storage/postgres"
	"4microservices/comment-service/storage/repo"
	"database/sql"
)

type IStorage interface {
	Comment() repo.CommentStorageI
}

type Pg struct {
	db          *sql.DB
	commentRepo repo.CommentStorageI
}

func NewStoragePg(db *sql.DB) *Pg {
	return &Pg{
		db:          db,
		commentRepo: postgres.NewCommentRepo(db),
	}
}

func (s Pg) Comment() repo.CommentStorageI {
	return s.commentRepo
}
