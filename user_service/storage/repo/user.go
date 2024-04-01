package repo

import (
	pb "4microservice/user-service/genproto/user_service"
	"github.com/golang/protobuf/ptypes/empty"
)

// UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.User, error)
	GetUser(id *pb.UserIdd) (*pb.GetUserResponse, error)
	UpdateUser(id *pb.User) (*pb.User, error)
	DeleteUser(id *pb.UserIdd) (*empty.Empty, error)
	GetAllUsers(req *pb.GetAllRequest) (*pb.GetAllResponse, error)
	CheckField(req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error)
	GetUserByEmail(req *pb.Email) (*pb.User, error)
}
