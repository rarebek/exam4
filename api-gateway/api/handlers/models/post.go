package models

type UserWithPosts struct {
	Posts []*Post `protobuf:"bytes,2,rep,name=posts,proto3" json:"posts,omitempty"`
}

type Post struct {
	Id       string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId   string      `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content  string      `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	ImageUrl string      `protobuf:"bytes,4,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	Title    string      `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Likes    int64       `protobuf:"varint,6,opt,name=likes,proto3" json:"likes,omitempty"`
	Dislikes int64       `protobuf:"varint,7,opt,name=dislikes,proto3" json:"dislikes,omitempty"`
	Views    int64       `protobuf:"varint,8,opt,name=views,proto3" json:"views,omitempty"`
	Comments []*Commentt `protobuf:"bytes,9,rep,name=comments,proto3" json:"comments,omitempty"`
}

type PostSwag struct {
	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId   string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Content  string `protobuf:"bytes,3,opt,name=content,proto3" json:"content,omitempty"`
	ImageUrl string `protobuf:"bytes,4,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
	Title    string `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Likes    int64  `protobuf:"varint,6,opt,name=likes,proto3" json:"likes,omitempty"`
	Dislikes int64  `protobuf:"varint,7,opt,name=dislikes,proto3" json:"dislikes,omitempty"`
	Views    int64  `protobuf:"varint,8,opt,name=views,proto3" json:"views,omitempty"`
}
