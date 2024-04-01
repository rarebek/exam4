package service

import (
	pb "4microservice/post-service/genproto/post_service"
	"4microservice/post-service/pkg/logger"
	grpcclient "4microservice/post-service/service/grpc_client"
	"4microservice/post-service/storage"
	"database/sql"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// PostService p
type PostService struct {
	storage storage.IStorage
	logger  logger.Logger
	client  grpcclient.IServiceManager
	pb.UnimplementedPostServiceServer
}

// NewPostService n
func NewPostService(db *sql.DB, log logger.Logger, client grpcclient.IServiceManager) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

// rpc CreatePost(Post) returns (Post);
// rpc UpdatePost(Post) returns (Post);
// rpc GetPost(PostID) returns (Post);
// rpc GetAllPosts(PostsRequest) returns (PostsResponse);
// rpc DeletePost(PostID) returns (Post);

func (p *PostService) CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	return p.storage.Post().CreatePost(req)
}

func (p *PostService) UpdatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	return p.storage.Post().UpdatePost(req)
}

func (p *PostService) DeletePost(ctx context.Context, req *pb.PostID) (*empty.Empty, error) {
	return p.storage.Post().DeletePost(req)
}

func (p *PostService) GetAllPosts(ctx context.Context, req *pb.PostsRequest) (*pb.UserWithPosts, error) {
	return p.storage.Post().GetAllPosts(req)
}

func (p *PostService) GetPost(ctx context.Context, req *pb.PostID) (*pb.Post, error) {
	return p.storage.Post().GetPost(req)
}

func (p *PostService) GetUserPostsByUserId(ctx context.Context, req *pb.UserId) (*pb.UserWithPosts, error) {
	return p.storage.Post().GetUserPostsByUserId(req)
}
