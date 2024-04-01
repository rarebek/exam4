package main

import (
	"4microservice/post-service/config"
	pb "4microservice/post-service/genproto/post_service"
	"4microservice/post-service/pkg/db"
	"4microservice/post-service/pkg/logger"
	"4microservice/post-service/service"
	grpcCLient "4microservice/post-service/service/grpc_client"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "post-service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {

		}
	}(log)

	log.Info("main: sqlConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	//connDBPsql, _, err := db.ConnectToDB(cfg)
	//if err != nil {
	//	log.Fatal("sql connection to postgres error", logger.Error(err))
	//}

	Conndb, _, err := db.ConnectToDB(cfg)

	client, err := grpcCLient.New(cfg)
	if err != nil {
		log.Fatal("error while adding grpc client", logger.Error(err))
	}

	postService := service.NewPostService(Conndb, log, client)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("failed to listen to: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	pb.RegisterPostServiceServer(s, postService)
	log.Info("main: server is running",
		logger.String("port", cfg.RPCPort))

	if err := s.Serve(lis); err != nil {
		log.Fatal("error while listening: %v", logger.Error(err))
	}
}
