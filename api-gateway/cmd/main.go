package main

import (
	"4microservice/api_gateway/api"
	"4microservice/api_gateway/config"
	"4microservice/api_gateway/pkg/logger"
	"4microservice/api_gateway/services"
	reds "4microservice/api_gateway/storage/redis"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/nats-io/nats.go"
)

type User struct {
	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username     string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email        string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password     string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	FirstName    string `protobuf:"bytes,5,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName     string `protobuf:"bytes,6,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Bio          string `protobuf:"bytes,7,opt,name=bio,proto3" json:"bio,omitempty"`
	Website      string `protobuf:"bytes,8,opt,name=website,proto3" json:"website,omitempty"`
	RefreshToken string `protobuf:"bytes,9,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
}

// @title API GATEWAY
// @version 1.0
// @description API GATEWAY SERVICE

// @host localhost:9090
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "api_gateway")
	pool := redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", "localhost", 6379))
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Error(err.Error())
	}
	defer nc.Close()

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error", logger.Error(err))
	}

	server := api.New(api.Option{
		Cfg:            cfg,
		Logger:         log,
		ServiceManager: serviceManager,
		Reds:           reds.NewRedisRepo(&pool),
		NatsCon:        nc,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("cannot run http server", logger.Error(err))
		panic(err)
	}

	// _, err = nc.Subscribe("user.created", func(msg *nats.Msg) {
	// 	handleUserCreatedMessage(msg)
	// })
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// select {}
}

func handleUserCreatedMessage(msg *nats.Msg) {
	var user User
	err := json.Unmarshal(msg.Data, &user)
	if err != nil {
		log.Println("Error decoding message:", err)
		return
	}
	log.Println("Received user created message:", user)
}
