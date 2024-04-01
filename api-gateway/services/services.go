package services

import (
	"4microservice/api_gateway/config"
	pbc "4microservice/api_gateway/genproto/comment_service"
	pbp "4microservice/api_gateway/genproto/post_service"
	pbu "4microservice/api_gateway/genproto/user_service"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() pbu.UserServiceClient
	PostService() pbp.PostServiceClient
	CommentService() pbc.CommentServiceClient
}

type serviceManager struct {
	userService    pbu.UserServiceClient
	postService    pbp.PostServiceClient
	commentService pbc.CommentServiceClient
}

func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}

func (s *serviceManager) PostService() pbp.PostServiceClient {
	return s.postService
}

func (s *serviceManager) CommentService() pbc.CommentServiceClient {
	return s.commentService
}

func NewServiceManager(cfg *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.PostServiceHost, cfg.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connComment, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.CommentServiceHost, cfg.CommentServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	serviceManager := &serviceManager{
		userService:    pbu.NewUserServiceClient(connUser),
		postService:    pbp.NewPostServiceClient(connPost),
		commentService: pbc.NewCommentServiceClient(connComment),
	}
	return serviceManager, nil
}
