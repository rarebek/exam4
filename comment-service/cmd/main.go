package main

import (
	"4microservices/comment-service/config"
	pb "4microservices/comment-service/genproto/comment_service"
	"4microservices/comment-service/pkg/db"
	"4microservices/comment-service/pkg/logger"
	"4microservices/comment-service/service"
	"net"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.LogLevel, "comment-service")
	defer func(l logger.Logger) {
		err := logger.Cleanup(l)
		if err != nil {
			log.Fatal(err.Error())
		}
	}(log)

	log.Info("main: sqlConfig",
		logger.String("host", cfg.PostgresHost),
		logger.Int("port", cfg.PostgresPort),
		logger.String("database", cfg.PostgresDatabase))

	connDB, _, err := db.ConnectDB(cfg)
	if err != nil {
		log.Fatal("sql connection error", logger.Error(err))
	}

	commentService := service.NewCommentService(connDB, log)

	lis, err := net.Listen("tcp", cfg.RPCPort)
	if err != nil {
		log.Fatal("cannot listen", logger.Error(err))
	}

	server := grpc.NewServer()
	pb.RegisterCommentServiceServer(server, commentService)
	log.Info("main: server running",
		logger.String("port", cfg.RPCPort))

	if err := server.Serve(lis); err != nil {
		log.Fatal("server cannot serve", logger.Error(err))
	}
}
