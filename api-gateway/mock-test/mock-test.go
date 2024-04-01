package mock_data

import (
	pbc "4microservice/api_gateway/genproto/comment_service"
	pbp "4microservice/api_gateway/genproto/post_service"
	pbu "4microservice/api_gateway/genproto/user_service"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"time"
)

type MockServiceClient interface {
	// CreateUser Create User
	CreateUser(*pbu.User) (*pbu.User, error)
	UpdateUser(id *pbu.User) (*pbu.User, error)
	GetAllUsers(req *pbu.GetAllRequest) (*pbu.GetAllResponse, error)
	CheckField(req *pbu.CheckFieldRequest) (*pbu.CheckFieldResponse, error)
	GetUserByEmail(req *pbu.Email) (*pbu.User, error)

	// CreatePost Post
	CreatePost(req *pbp.Post) (*pbp.Post, error)
	UpdatePost(req *pbp.Post) (*pbp.Post, error)
	DeletePost(req *pbp.PostID) (*empty.Empty, error)
	GetPost(req *pbp.PostID) (*pbp.Post, error)
	GetAllPosts(req *pbp.PostsRequest) (*pbp.UserWithPosts, error)
	GetUserPostsByUserId(req *pbp.UserId) (*pbp.UserWithPosts, error)

	// CreateComment Comment
	CreateComment(req *pbc.Comment) (*pbc.Comment, error)
	UpdateComment(req *pbc.Comment) (*pbc.Comment, error)
	GetComment(req *pbc.CommentId) (*pbc.Comment, error)
	GetAllComments(req *pbc.CommentsRequest) (*pbc.CommentsResponse, error)
	DeleteComment(req *pbc.CommentId) (*pbc.Comment, error)
}

type mockServiceClient struct {
}

func NewMockServiceClient() MockServiceClient {
	return &mockServiceClient{}
}
func (r *mockServiceClient) CheckField(request *pbu.CheckFieldRequest) (*pbu.CheckFieldResponse, error) {
	return &pbu.CheckFieldResponse{Unique: true}, nil
}

func (r *mockServiceClient) CreateUser(user *pbu.User) (*pbu.User, error) {

	return &pbu.User{
		Id:           user.Id,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		Email:        user.Email,
		Password:     user.Password,
		Bio:          user.Bio,
		Website:      user.Website,
		RefreshToken: "mock token",
	}, nil
}

func (r *mockServiceClient) GetUserByID(in *pbu.UserIdd) (*pbu.User, error) {

	return &pbu.User{
		Id:           in.Id,
		FirstName:    "Nodirbek",
		LastName:     "Nomonov",
		Username:     "Mock username",
		Bio:          "Mock bio",
		Email:        "Mock Email",
		Password:     "MockPassword",
		Website:      "Mock website",
		RefreshToken: "MOck token",
	}, nil
}

func (r *mockServiceClient) GetUserByEmail(in *pbu.Email) (*pbu.User, error) {

	return &pbu.User{
		FirstName:    "Nodirbek",
		LastName:     "Nomonov",
		Username:     "Mock username",
		Bio:          "Mock bio",
		Email:        "Mock Email",
		Password:     "MockPassword",
		Website:      "Mock website",
		RefreshToken: "MOck token",
	}, nil
}

func (r *mockServiceClient) GetAllUsers(in *pbu.GetAllRequest) (*pbu.GetAllResponse, error) {
	results := []*pbu.GetUserResponse{
		{
			Id:        "b124c9a5-3ae3-4597-9540-c9c2b06ed050",
			FirstName: "Nodirbek",
			LastName:  "Nomonov",
			Username:  "Mock username",
			Bio:       "Mock bio",
			Email:     "Mock Email",
			Password:  "Mock password",
			Website:   "Mock site",
		},
		{
			Id:        "b124c9a5-3ae3-4597-9540-c9c2b06ed050",
			FirstName: "Nodirbek2",
			LastName:  "Nomonov2",
			Username:  "Mock username",
			Bio:       "Mock bio",
			Email:     "Mock Email",
			Password:  "Mock password",
			Website:   "Mock site",
		},
	}

	return &pbu.GetAllResponse{
		Users: results,
	}, nil
}

func (r *mockServiceClient) UpdateUser(res *pbu.User) (*pbu.User, error) {

	return &pbu.User{
		Id:        res.Id,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Username:  res.Username,
		Email:     res.Email,
		Password:  res.Password,
		Bio:       res.Bio,
		Website:   res.Website,
	}, nil

}

func (r *mockServiceClient) DeleteUserByID(in *pbu.UserIdd) (*empty.Empty, error) {

	return nil, nil

}

func (r *mockServiceClient) CheckUniqueRequest(in *pbu.CheckFieldRequest) (*pbu.CheckFieldResponse, error) {

	return &pbu.CheckFieldResponse{
		Unique: true,
	}, nil
}

// CreatePost Post
func (r *mockServiceClient) CreatePost(post *pbp.Post) (*pbp.Post, error) {
	if post.Id == "" {
		post.Id = uuid.New().String()
	}

	nowT := time.Now()

	return &pbp.Post{
		Id:        post.Id,
		UserId:    post.UserId,
		Content:   post.Content,
		Title:     post.Title,
		Likes:     post.Likes,
		Dislikes:  post.Dislikes,
		Views:     post.Views,
		CreatedAt: nowT.String(),
		UpdatedAt: nowT.String(),
	}, nil

}

func (r *mockServiceClient) GetPost(post *pbp.PostID) (*pbp.Post, error) {
	if post.Id == "" {
		post.Id = uuid.New().String()
	}

	nowT := time.Now()

	return &pbp.Post{
		Id:        post.Id,
		Title:     "Mock title",
		CreatedAt: nowT.String(),
		UpdatedAt: nowT.String(),
	}, nil

}

func (r *mockServiceClient) GetUserPostsByUserId(req *pbp.UserId) (*pbp.UserWithPosts, error) {
	posts := []*pbp.Post{
		{
			Id:        "2926aff0-dd15-49f0-9a55-265554c4628e",
			UserId:    req.UserId,
			Content:   "mock content 1",
			Title:     "mock title 1",
			CreatedAt: "2000-03-04",
			UpdatedAt: "2000-03-04",
		},
		{
			Id:        "4a9c939e-3a5d-404b-8585-bda6be573e3f",
			UserId:    req.UserId,
			Content:   "mock content 2",
			Title:     "mock title 2",
			CreatedAt: "1956-03-04",
			UpdatedAt: "1999-03-04",
		},
	}
	return &pbp.UserWithPosts{
		Posts: posts,
	}, nil
}

func (r *mockServiceClient) GetAllPosts(req *pbp.PostsRequest) (*pbp.UserWithPosts, error) {
	posts := []*pbp.Post{
		{
			Id:        "2926aff0-dd15-49f0-9a55-265554c4628e",
			Content:   "mock content 1",
			Title:     "mock title 1",
			CreatedAt: "2000-03-04",
			UpdatedAt: "2000-03-04",
		},
		{
			Id:        "4a9c939e-3a5d-404b-8585-bda6be573e3f",
			Content:   "mock content 2",
			Title:     "mock title 2",
			CreatedAt: "1956-03-04",
			UpdatedAt: "1999-03-04",
		},
	}
	return &pbp.UserWithPosts{
		Posts: posts,
	}, nil
}

func (r *mockServiceClient) UpdatePost(res *pbp.Post) (*pbp.Post, error) {

	return &pbp.Post{
		Id:        res.Id,
		UserId:    res.UserId,
		Content:   res.Content,
		Title:     res.Title,
		Likes:     res.Likes,
		Dislikes:  res.Dislikes,
		Views:     res.Views,
		CreatedAt: "2000-03-04",
		UpdatedAt: time.Now().String(),
	}, nil
}

func (r *mockServiceClient) DeletePost(*pbp.PostID) (*empty.Empty, error) {
	return nil, nil
}

// CreateComment Comment
func (r *mockServiceClient) CreateComment(res *pbc.Comment) (*pbc.Comment, error) {
	if res.Id == "" {
		res.Id = uuid.New().String()
	}
	nowT := time.Now().String()
	return &pbc.Comment{
		Id:        res.Id,
		Content:   res.Content,
		UserId:    res.UserId,
		PostId:    res.PostId,
		CreatedAt: nowT,
		UpdatedAt: nowT,
	}, nil

}

func (r *mockServiceClient) GetComment(res *pbc.CommentId) (*pbc.Comment, error) {
	if res.Id == "" {
		res.Id = uuid.New().String()
	}
	nowT := time.Now().String()
	return &pbc.Comment{
		Id:        res.Id,
		CreatedAt: nowT,
		UpdatedAt: nowT,
	}, nil

}

func (r *mockServiceClient) GetAllCommentsByPostID(req *pbc.CommentsRequest) (*pbc.CommentsResponse, error) {
	comments := []*pbc.Comment{
		{
			Id:        "db884df9-eaaa-4a72-b829-7a83c0faed95",
			Content:   "mock content 1",
			UserId:    "608b56be-dc0c-438e-98a9-d7804546fa39",
			PostId:    req.PostId,
			CreatedAt: "2010-08-02",
			UpdatedAt: "2010-08-02",
		},
		{
			Id:        "4dccfbac-a2eb-43ca-a89a-0be4341f3c21",
			Content:   "mock content 2",
			UserId:    "30c6e251-3c4a-462a-a7f6-09edd9a9f902",
			PostId:    req.PostId,
			CreatedAt: "2010-08-02",
			UpdatedAt: "2010-08-02",
		},
	}

	return &pbc.CommentsResponse{
		Comments: comments,
	}, nil
}

func (r *mockServiceClient) GetAllComments(req *pbc.CommentsRequest) (*pbc.CommentsResponse, error) {
	comments := []*pbc.Comment{
		{
			Id:        "db884df9-eaaa-4a72-b829-7a83c0faed95",
			Content:   "mock content 1",
			UserId:    "608b56be-dc0c-438e-98a9-d7804546fa39",
			PostId:    req.PostId,
			CreatedAt: "2010-08-02",
			UpdatedAt: "2010-08-02",
		},
		{
			Id:        "4dccfbac-a2eb-43ca-a89a-0be4341f3c21",
			Content:   "mock content 2",
			UserId:    "30c6e251-3c4a-462a-a7f6-09edd9a9f902",
			PostId:    req.PostId,
			CreatedAt: "2010-08-02",
			UpdatedAt: "2010-08-02",
		},
	}

	return &pbc.CommentsResponse{
		Comments: comments,
	}, nil
}

func (r *mockServiceClient) UpdateComment(res *pbc.Comment) (*pbc.Comment, error) {

	return &pbc.Comment{
		Id:        res.Id,
		Content:   res.Content,
		UserId:    res.UserId,
		PostId:    res.PostId,
		CreatedAt: "2010-08-02",
		UpdatedAt: time.Now().String(),
	}, nil

}

func (r *mockServiceClient) DeleteComment(req *pbc.CommentId) (*pbc.Comment, error) {
	return &pbc.Comment{
		Content: "Mock content",
	}, nil
}
