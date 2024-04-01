package service

import (
	pbc "4microservice/user-service/genproto/comment_service"
	pbp "4microservice/user-service/genproto/post_service"
	pb "4microservice/user-service/genproto/user_service"
	l "4microservice/user-service/pkg/logger"
	grpcClient "4microservice/user-service/service/grpc_client"
	"4microservice/user-service/storage"
	"context"
	"database/sql"

	"github.com/golang/protobuf/ptypes/empty"
)

// UserService ...
type UserService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcClient.IServiceManager
	pb.UnimplementedUserServiceServer
}

// NewUserService ...
func NewUserService(db *sql.DB, log l.Logger, client grpcClient.IServiceManager) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	return s.storage.User().CreateUser(req)
}

func (s *UserService) GetUser(ctx context.Context, req *pb.UserIdd) (*pb.GetUserResponse, error) {
	var (
		response = pb.GetUserResponse{Posts: []*pb.Postt{}}
	)
	user, err := s.storage.User().GetUser(req)
	if err != nil {
		return nil, err
	}
	response.Id = user.Id
	response.Website = user.Website
	response.Bio = user.Bio
	response.Email = user.Email
	response.Password = user.Password
	response.LastName = user.LastName
	response.FirstName = user.FirstName
	response.Username = user.Username
	posts, err := s.client.PostService().GetUserPostsByUserId(ctx, &pbp.UserId{UserId: user.Id})
	if err != nil {
		return nil, err
	}
	var (
		userPosts []*pb.Postt
	)
	for _, post := range posts.Posts {
		var postComments []*pb.Commentt
		comments, err := s.client.CommentService().GetAllComments(ctx, &pbc.CommentsRequest{PostId: post.Id})
		if err != nil {
			return nil, err
		}
		for _, comment := range comments.Comments {
			var postComment pb.Commentt
			postComment.Id = comment.Id
			postComment.PostId = comment.PostId
			postComment.Content = comment.Content
			postComment.UpdatedAt = comment.UpdatedAt
			postComment.CreatedAt = comment.CreatedAt
			postComment.UserId = comment.UserId

			postComments = append(postComments, &postComment)
		}
		var userPost pb.Postt
		userPost.Id = post.Id
		userPost.UserId = post.UserId
		userPost.Content = post.Content
		userPost.Views = post.Views
		userPost.Likes = post.Likes
		userPost.Dislikes = post.Dislikes
		userPost.Title = post.Title
		userPost.ImageUrl = post.ImageUrl
		userPost.Comments = postComments

		userPosts = append(userPosts, &userPost)
	}
	response.Posts = userPosts

	return &response, nil
}

func (s *UserService) GetAllUser(ctx context.Context, req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	users, err := s.storage.User().GetAllUsers(req)
	if err != nil {
		return nil, err
	}

	for _, user := range users.Users {
		posts, err := s.client.PostService().GetUserPostsByUserId(ctx, &pbp.UserId{UserId: user.Id})
		if err != nil {
			return nil, err
		}
		var (
			userPosts []*pb.Postt
		)
		for _, post := range posts.Posts {
			var postComments []*pb.Commentt
			comments, err := s.client.CommentService().GetAllComments(ctx, &pbc.CommentsRequest{PostId: post.Id})
			if err != nil {
				return nil, err
			}
			for _, comment := range comments.Comments {
				var postComment pb.Commentt
				postComment.Id = comment.Id
				postComment.PostId = comment.PostId
				postComment.Content = comment.Content
				postComment.UpdatedAt = comment.UpdatedAt
				postComment.CreatedAt = comment.CreatedAt
				postComment.UserId = comment.UserId

				postComments = append(postComments, &postComment)
			}
			var userPost pb.Postt
			userPost.Id = post.Id
			userPost.UserId = post.UserId
			userPost.Content = post.Content
			userPost.Views = post.Views
			userPost.Likes = post.Likes
			userPost.Dislikes = post.Dislikes
			userPost.Title = post.Title
			userPost.ImageUrl = post.ImageUrl
			userPost.Comments = postComments

			userPosts = append(userPosts, &userPost)
		}
		user.Posts = userPosts
	}
	return users, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.UserIdd) (*empty.Empty, error) {
	return s.storage.User().DeleteUser(req)
}

func (s *UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	return s.storage.User().UpdateUser(req)
}

func (s *UserService) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	return s.storage.User().CheckField(req)
}

func (s *UserService) GetUserByEmail(ctx context.Context, req *pb.Email) (*pb.User, error) {
	return s.storage.User().GetUserByEmail(req)
}
