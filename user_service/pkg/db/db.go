package db

import (
	"4microservice/user-service/config"
	"context"
	"database/sql"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" //postgres drivers
	"github.com/nats-io/nats.go"
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

func ConnectToRedis(cfg config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	return client
}

func ConnectToMongoDB(uri, dbName string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func ConnectToNATS() (*nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	return nc, nil
}
