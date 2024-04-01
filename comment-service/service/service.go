package service

import (
	pbc "4microservices/comment-service/genproto/comment_service"
	l "4microservices/comment-service/pkg/logger"
	"4microservices/comment-service/storage"
	"context"
	"database/sql"
)

type CommentService struct {
	storage storage.IStorage
	logger  l.Logger
	pbc.UnimplementedCommentServiceServer
}

func NewCommentService(db *sql.DB, log l.Logger) *CommentService {
	return &CommentService{
		storage: storage.NewStoragePg(db),
		logger:  log,
	}
}

//rpc CreateComment(Comment) returns (Comment);
//rpc UpdateComment(Comment) returns (Comment);
//rpc GetComment(CommentId) returns (Comment);
//rpc GetAllComments(CommentsRequest) returns (CommentsResponse);
//rpc DeleteComment(CommentId) returns (Comment);

func (c *CommentService) CreateComment(ctx context.Context, req *pbc.Comment) (*pbc.Comment, error) {
	return c.storage.Comment().CreateComment(req)
}

func (c *CommentService) UpdateComment(ctx context.Context, req *pbc.Comment) (*pbc.Comment, error) {
	return c.storage.Comment().UpdateComment(req)
}

func (c *CommentService) GetComment(ctx context.Context, req *pbc.CommentId) (*pbc.Comment, error) {
	return c.storage.Comment().GetComment(req)
}

func (c *CommentService) GetAllComments(ctx context.Context, req *pbc.CommentsRequest) (*pbc.CommentsResponse, error) {
	return c.storage.Comment().GetAllComments(req)
}

func (c *CommentService) DeleteComment(cyx context.Context, req *pbc.CommentId) (*pbc.Comment, error) {
	return c.storage.Comment().DeleteComment(req)
}
