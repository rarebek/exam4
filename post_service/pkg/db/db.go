package db

import (
	"4microservice/post-service/config"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" //postgres drivers
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func ConnectToDB(cfg config.Config) (*sql.DB, func(), error) {
	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	connDb, err := sql.Open("postgres", psqlString)
	if err != nil {
		return nil, nil, err
	}

	cleanUpFunc := func() {
		connDb.Close()
	}
	return connDb, cleanUpFunc, nil
}

func New(cfg config.Config) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	return client, nil
}
