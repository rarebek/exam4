package models

import (
	"google.golang.org/protobuf/runtime/protoimpl"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v3"
	"github.com/go-ozzo/ozzo-validation/v3/is"
)

type RegisterUserResponse struct {
	Message string
}
type User struct {
	Id           string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username     string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email        string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password     string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	FirstName    string `protobuf:"bytes,5,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName     string `protobuf:"bytes,6,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Bio          string `protobuf:"bytes,7,opt,name=bio,proto3" json:"bio,omitempty"`
	Website      string `protobuf:"bytes,8,opt,name=website,proto3" json:"website,omitempty"`
	RefreshToken string `protobuf:"bytes,9,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
}

type UserModelForRegister struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username  string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email     string `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password  string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	FirstName string `protobuf:"bytes,5,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName  string `protobuf:"bytes,6,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Bio       string `protobuf:"bytes,7,opt,name=bio,proto3" json:"bio,omitempty"`
	Website   string `protobuf:"bytes,8,opt,name=website,proto3" json:"website,omitempty"`
	OTP       string `protobuf:"bytes,9,opt,name=otp,proto3" json:"otp,omitempty"`
}

type GetAllResponse struct {
	Users []*GetUserResponse `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"`
}

type GetUserResponse struct {
	Id           string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username     string   `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	Email        string   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password     string   `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	FirstName    string   `protobuf:"bytes,5,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName     string   `protobuf:"bytes,6,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	Bio          string   `protobuf:"bytes,7,opt,name=bio,proto3" json:"bio,omitempty"`
	Website      string   `protobuf:"bytes,8,opt,name=website,proto3" json:"website,omitempty"`
	RefreshToken string   `protobuf:"bytes,9,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	Posts        []*Postt `protobuf:"bytes,10,rep,name=posts,proto3" json:"posts,omitempty"`
}

type Postt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Id            string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        string      `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content       string      `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	ImageUrl      string      `protobuf:"bytes,4,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	Title         string      `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Likes         int64       `protobuf:"varint,6,opt,name=likes,proto3" json:"likes,omitempty"`
	Dislikes      int64       `protobuf:"varint,7,opt,name=dislikes,proto3" json:"dislikes,omitempty"`
	Views         int64       `protobuf:"varint,8,opt,name=views,proto3" json:"views,omitempty"`
	Comments      []*Commentt `protobuf:"bytes,9,rep,name=comments,proto3" json:"comments,omitempty"`
}

type Commentt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId        string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PostId        string `protobuf:"bytes,3,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	Content       string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	CreatedAt     string `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     string `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt     string `protobuf:"bytes,7,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
}
type RefreshTokenUpdate struct {
	RefreshToken string `json:"refreshToken"`
}

type ReturningForFrontEnd struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Age         int64  `json:"age"`
	AccessToken string `json:"accessToken"`
}

type RegisterModelResponse struct {
	Result string `json:"result"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}

type VerifyRequest struct {
	Code string
	User User
}

type VerifyResponse struct {
	Success bool
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 15), validation.Match(regexp.MustCompile("[a-z]|[A-Z][0-9]"))),
		validation.Field(&u.Username, validation.Required, validation.Length(3, 50), validation.Match(regexp.MustCompile("^[A-Z][a-z]*$"))),
	)
}
