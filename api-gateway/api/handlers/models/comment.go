package models

type Comment struct {
	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId    string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	PostId    string `protobuf:"bytes,3,opt,name=post_id,json=postId,proto3" json:"post_id,omitempty"`
	Content   string `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	CreatedAt string `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt string `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	DeletedAt string `protobuf:"bytes,7,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty"`
}

type CommentsResponse struct {
	Comments []*Comment `protobuf:"bytes,1,rep,name=comments,proto3" json:"comments,omitempty"`
	Count    int64      `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}
