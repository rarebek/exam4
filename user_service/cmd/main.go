package main

import (
	"4microservice/user-service/config"
	pb "4microservice/user-service/genproto/user_service"
	"4microservice/user-service/pkg/db"
	"4microservice/user-service/pkg/logger"
	"4microservice/user-service/service"
	grpcClient "4microservice/user-service/service/grpc_client"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "user-service")

	log.Info("main: sqlConfig",
		logger.String("address", cfg.PostgresHost),
		logger.String("database", cfg.PostgresDatabase))

	connDB, _, err := db.ConnectToDB(cfg)
	if err != nil {
		log.Fatal("sql connection to postgres error", logger.Error(err))
	}

	// conn, err := db.ConnectToNATS()
	// if err != nil {
	// 	log.Fatal("error while nats connection", logger.Error(err))
	// }

	client, err := grpcClient.New(cfg)
	if err != nil {
		log.Fatal("error while creating new client", logger.Error(err))
	}

	userService := service.NewUserService(connDB, log, client)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, userService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("Error while listening: %v", logger.Error(err))
	}
}
