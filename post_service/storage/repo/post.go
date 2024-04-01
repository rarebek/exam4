package repo

import (
	pb "4microservice/post-service/genproto/post_service"
	"github.com/golang/protobuf/ptypes/empty"
)

type PostStorageI interface {
	CreatePost(req *pb.Post) (*pb.Post, error)
	UpdatePost(req *pb.Post) (*pb.Post, error)
	DeletePost(req *pb.PostID) (*empty.Empty, error)
	GetPost(req *pb.PostID) (*pb.Post, error)
	GetAllPosts(req *pb.PostsRequest) (*pb.UserWithPosts, error)
	GetUserPostsByUserId(req *pb.UserId) (*pb.UserWithPosts, error)
}
