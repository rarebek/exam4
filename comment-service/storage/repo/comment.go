package repo

import pbc "4microservices/comment-service/genproto/comment_service"

//rpc CreateComment(Comment) returns (Comment);
//rpc UpdateComment(Comment) returns (Comment);
//rpc GetComment(CommentId) returns (Comment);
//rpc GetAllComments(CommentsRequest) returns (CommentsResponse);
//rpc DeleteComment(CommentId) returns (Comment);

type CommentStorageI interface {
	CreateComment(req *pbc.Comment) (*pbc.Comment, error)
	UpdateComment(req *pbc.Comment) (*pbc.Comment, error)
	GetComment(req *pbc.CommentId) (*pbc.Comment, error)
	GetAllComments(req *pbc.CommentsRequest) (*pbc.CommentsResponse, error)
	DeleteComment(req *pbc.CommentId) (*pbc.Comment, error)
}
