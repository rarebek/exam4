package grpc_client

import (
	"4microservice/post-service/config"
	pbc "4microservice/post-service/genproto/comment_service"
	pbu "4microservice/post-service/genproto/user_service"
	"4microservice/post-service/pkg/logger"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

type IServiceManager interface {
	CommentService() pbc.CommentServiceClient
	UserService() pbu.UserServiceClient
}

type serviceManager struct {
	cfg            config.Config
	commentService pbc.CommentServiceClient
	userService    pbu.UserServiceClient
}

func New(cfg config.Config) (IServiceManager, error) {
	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal(""+
			"error while dialing to the post service", logger.Error(err))
	}

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while dialing to the user service", logger.Error(err))
	}

	return &serviceManager{
		cfg:            cfg,
		commentService: pbc.NewCommentServiceClient(connComment),
		userService:    pbu.NewUserServiceClient(connUser),
	}, nil
}

func (s *serviceManager) CommentService() pbc.CommentServiceClient {
	return s.commentService
}

func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}
