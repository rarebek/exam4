package grpcClient

import (
	"4microservice/user-service/config"
	pbc "4microservice/user-service/genproto/comment_service"
	pbp "4microservice/user-service/genproto/post_service"
	"4microservice/user-service/pkg/logger"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

type IServiceManager interface {
	PostService() pbp.PostServiceClient
	CommentService() pbc.CommentServiceClient
}

type serviceManager struct {
	cfg            config.Config
	postService    pbp.PostServiceClient
	commentService pbc.CommentServiceClient
}

func New(cfg config.Config) (IServiceManager, error) {
	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while dialing to the post service", logger.Error(err))
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while dialing to the post service", logger.Error(err))
	}

	return &serviceManager{
		cfg:            cfg,
		postService:    pbp.NewPostServiceClient(connPost),
		commentService: pbc.NewCommentServiceClient(connComment),
	}, nil
}

func (s *serviceManager) PostService() pbp.PostServiceClient {
	return s.postService
}

func (s *serviceManager) CommentService() pbc.CommentServiceClient {
	return s.commentService
}
