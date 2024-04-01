package storage

import (
	"4microservice/user-service/storage/postgres"
	"4microservice/user-service/storage/repo"
	"database/sql"

	"github.com/nats-io/nats.go"
)

// IStorage ...
type IStorage interface {
	User() repo.UserStorageI
}

type Pg struct {
	db       *sql.DB
	conn     *nats.Conn
	userRepo repo.UserStorageI
}

// NewStoragePg ...
func NewStoragePg(db *sql.DB) *Pg {
	return &Pg{
		db:       db,
		userRepo: postgres.NewUserRepo(db),
	}
}

func (s Pg) User() repo.UserStorageI {
	return s.userRepo
}
