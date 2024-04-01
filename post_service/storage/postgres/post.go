package postgres

import (
	pb "4microservice/post-service/genproto/post_service"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
)

type PostRepo struct {
	db *sql.DB
}

// NewPostRepo newpostrepo
func NewPostRepo(db *sql.DB) *PostRepo {
	return &PostRepo{db: db}
}

func (p *PostRepo) CreatePost(req *pb.Post) (*pb.Post, error) {
	req.Id = uuid.NewString()
	var updatedAT sql.NullTime
	query := `INSERT INTO posts(id, 
                                user_id, 
                  				content, 
                  				image_url, 
                  				title, 
                  				likes, 
                  				dislikes, 
                  				views) VALUES($1, $2, $3, $4, $5, $6, $7, $8)  
             RETURNING id, user_id, content, image_url, title, likes, dislikes, views, created_at, updated_at`

	var respPost pb.Post

	rowPost := p.db.QueryRow(
		query,
		req.Id,
		req.UserId,
		req.Content,
		req.ImageUrl,
		req.Title,
		req.Likes,
		req.Dislikes,
		req.Views,
	)

	if err := rowPost.Scan(
		&respPost.Id,
		&respPost.UserId,
		&respPost.Content,
		&respPost.ImageUrl,
		&respPost.Title,
		&respPost.Likes,
		&respPost.Dislikes,
		&respPost.Views,
		&respPost.CreatedAt,
		&updatedAT,
	); err != nil {
		return nil, err
	}

	if updatedAT.Valid {
		respPost.UpdatedAt = updatedAT.Time.String()
	}

	return &respPost, nil
}

func (p *PostRepo) GetPost(req *pb.PostID) (*pb.Post, error) {
	var updatedAT sql.NullTime
	query := `SELECT id, 
       				 user_id, 
       				 content, 
       				 image_url, 
       				 title, 
       				 likes, 
       				 dislikes, 
       				 views,
       				 created_at,
       				 updated_at
			  FROM posts 
			  WHERE id = $1`

	var respPost pb.Post
	rowPost := p.db.QueryRow(query, req.Id)

	if err := rowPost.Scan(&respPost.Id,
		&respPost.UserId,
		&respPost.Content,
		&respPost.ImageUrl,
		&respPost.Title,
		&respPost.Likes,
		&respPost.Dislikes,
		&respPost.Views,
		&respPost.CreatedAt,
		&updatedAT,
	); err != nil {
		return nil, err
	}

	if updatedAT.Valid {
		respPost.UpdatedAt = updatedAT.Time.String()
	}

	return &respPost, nil
}

func (p *PostRepo) UpdatePost(req *pb.Post) (*pb.Post, error) {
	var updatedAT sql.NullTime
	query := `UPDATE posts 
			  SET content = $1, 
			      image_url = $2, 
			      title = $3 
			  WHERE id = $4
			  RETURNING id, 
			  			user_id, 
			  			content, 
			  			image_url, 
			  			title, 
			  			likes, 
			  			dislikes, 
			  			views, 
			  			created_at, 
			  			updated_at`
	var respPost pb.Post
	rowPost := p.db.QueryRow(query, req.Content, req.ImageUrl, req.Title, req.Id)
	if err := rowPost.Scan(&respPost.Id,
		&respPost.UserId,
		&respPost.Content,
		&respPost.ImageUrl,
		&respPost.Title,
		&respPost.Likes,
		&respPost.Dislikes,
		&respPost.Views,
		&respPost.CreatedAt,
		&updatedAT,
	); err != nil {
		return nil, err
	}

	if updatedAT.Valid {
		respPost.UpdatedAt = updatedAT.Time.String()
	}

	return &respPost, nil
}

func (p *PostRepo) DeletePost(req *pb.PostID) (*empty.Empty, error) {
	query := `DELETE FROM posts WHERE id = $1`
	p.db.QueryRow(query, req.Id)
	return &empty.Empty{}, nil
}

func (p *PostRepo) GetAllPosts(req *pb.PostsRequest) (*pb.UserWithPosts, error) {
	offset := (req.Page - 1) * req.Limit
	query := `SELECT id,
						 user_id,
						 content,
						 image_url,
						 title,
						 likes,
						 dislikes,
						 views,
						 created_at,
						 updated_at
				  FROM posts LIMIT $1 OFFSET $2`
	rows, err := p.db.Query(query, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	var updatedAT sql.NullTime
	var posts pb.UserWithPosts
	for rows.Next() {
		var respPost pb.Post
		err := rows.Scan(&respPost.Id,
			&respPost.UserId,
			&respPost.Content,
			&respPost.ImageUrl,
			&respPost.Title,
			&respPost.Likes,
			&respPost.Dislikes,
			&respPost.Views,
			&respPost.CreatedAt,
			&updatedAT)
		if err != nil {
			return nil, err
		}
		if updatedAT.Valid {
			respPost.UpdatedAt = updatedAT.Time.String()
		}
		posts.Posts = append(posts.Posts, &respPost)
	}
	return &posts, nil
}

func (p *PostRepo) GetUserPostsByUserId(req *pb.UserId) (*pb.UserWithPosts, error) {
	query := `SELECT id, content, image_url, title, likes, dislikes, views, created_at, updated_at from posts WHERE user_id = $1`
	rows, err := p.db.Query(query, req.UserId)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	posts := pb.UserWithPosts{
		Posts: []*pb.Post{},
	}

	for rows.Next() {
		var UpdatedAT sql.NullTime
		var post pb.Post
		err := rows.Scan(&post.Id, &post.Content, &post.ImageUrl, &post.Title, &post.Likes, &post.Dislikes, &post.Views, &post.CreatedAt, &UpdatedAT)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		if UpdatedAT.Valid {
			post.UpdatedAt = UpdatedAT.Time.String()
		}
		newPost := pb.Post{
			Id:        post.Id,
			Content:   post.Content,
			ImageUrl:  post.ImageUrl,
			Title:     post.Title,
			Likes:     post.Likes,
			Dislikes:  post.Dislikes,
			Views:     post.Views,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}
		posts.Posts = append(posts.Posts, &newPost)
	}

	return &posts, nil
}
