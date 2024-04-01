package cassandra

import (
	pb "4microservice/post-service/genproto/post_service"
	"github.com/gocql/gocql"
	"github.com/golang/protobuf/ptypes/empty"
)

type PostRepo struct {
	session *gocql.Session
}

func NewPostRepo(clusterHosts []string, keyspace string) (*PostRepo, error) {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &PostRepo{session: session}, nil
}

func (p *PostRepo) Close() {
	p.session.Close()
}

func (p *PostRepo) CreatePost(req *pb.Post) (*pb.Post, error) {
	req.Id = gocql.TimeUUID().String()
	query := `INSERT INTO posts (id, user_id, content, image_url, title, likes, dislikes, views) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if err := p.session.Query(query,
		req.Id,
		req.UserId,
		req.Content,
		req.ImageUrl,
		req.Title,
		req.Likes,
		req.Dislikes,
		req.Views).Exec(); err != nil {
		return nil, err
	}
	return req, nil
}

func (p *PostRepo) GetPost(req *pb.PostID) (*pb.Post, error) {
	var post pb.Post
	query := `SELECT id, user_id, content, image_url, title, likes, dislikes, views FROM posts WHERE id = ?`
	err := p.session.Query(query, req.Id).Scan(&post.Id, &post.UserId, &post.Content, &post.ImageUrl, &post.Title, &post.Likes, &post.Dislikes, &post.Views)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p *PostRepo) UpdatePost(req *pb.Post) (*pb.Post, error) {
	query := `UPDATE posts SET content = ?, image_url = ?, title = ? WHERE id = ?`
	if err := p.session.Query(query, req.Content, req.ImageUrl, req.Title, req.Id).Exec(); err != nil {
		return nil, err
	}
	return req, nil
}

func (p *PostRepo) DeletePost(req *pb.PostID) (*empty.Empty, error) {
	query := `DELETE FROM posts WHERE id = ?`
	if err := p.session.Query(query, req.Id).Exec(); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (p *PostRepo) GetAllPosts(req *pb.PostsRequest) (*pb.UserWithPosts, error) {
	var posts pb.UserWithPosts
	query := `SELECT id, user_id, content, image_url, title, likes, dislikes, views FROM posts LIMIT ? OFFSET ?`
	iter := p.session.Query(query, req.Limit, (req.Page-1)*req.Limit).Iter()
	var post pb.Post
	for iter.Scan(&post.Id, &post.UserId, &post.Content, &post.ImageUrl, &post.Title, &post.Likes, &post.Dislikes, &post.Views) {
		posts.Posts = append(posts.Posts, &post)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return &posts, nil
}

func (p *PostRepo) GetUserPostsByUserId(req *pb.UserId) (*pb.UserWithPosts, error) {
	var posts pb.UserWithPosts
	query := `SELECT id, content, image_url, title, likes, dislikes, views FROM posts WHERE user_id = ?`
	iter := p.session.Query(query, req.UserId).Iter()
	var post pb.Post
	for iter.Scan(&post.Id, &post.Content, &post.ImageUrl, &post.Title, &post.Likes, &post.Dislikes, &post.Views) {
		posts.Posts = append(posts.Posts, &pb.Post{
			Id:        post.Id,
			UserId:    req.UserId,
			Content:   post.Content,
			ImageUrl:  post.ImageUrl,
			Title:     post.Title,
			Likes:     post.Likes,
			Dislikes:  post.Dislikes,
			Views:     post.Views,
			CreatedAt: "",
			UpdatedAt: "",
		})
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return &posts, nil
}
