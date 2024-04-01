package postgres

import (
	pbc "4microservices/comment-service/genproto/comment_service"
	"database/sql"
)

type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

//rpc CreateComment(Comment) returns (Comment);
//rpc UpdateComment(Comment) returns (Comment);
//rpc GetComment(CommentId) returns (Comment);
//rpc GetAllComments(CommentsRequest) returns (CommentsResponse);
//rpc DeleteComment(CommentId) returns (Comment);

func (c *CommentRepo) CreateComment(req *pbc.Comment) (*pbc.Comment, error) {
	var updatedAT sql.NullTime
	query := `INSERT INTO comments(id, user_id, post_id, content) VALUES($1, $2, $3, $4)  RETURNING id, user_id, post_id, content, created_at, updated_at`

	var respComment pbc.Comment

	rowPost := c.db.QueryRow(
		query,
		req.Id,
		req.UserId,
		req.PostId,
		req.Content,
	)

	if err := rowPost.Scan(
		&respComment.Id,
		&respComment.UserId,
		&respComment.PostId,
		&respComment.Content,
		&respComment.CreatedAt,
		&updatedAT,
	); err != nil {
		return nil, err
	}

	if updatedAT.Valid {
		respComment.UpdatedAt = updatedAT.Time.String()
	}

	return &respComment, nil
}

func (c *CommentRepo) GetComment(req *pbc.CommentId) (*pbc.Comment, error) {
	var updatedAT sql.NullTime
	query := `SELECT id, 
       				 user_id, 
       				 post_id,
       				 content, 
       				 created_at,
       				 updated_at
			  FROM comments 
			  WHERE id = $1`

	var respComment pbc.Comment
	rowPost := c.db.QueryRow(query, req.Id)

	if err := rowPost.Scan(&respComment.Id,
		&respComment.UserId,
		&respComment.PostId,
		&respComment.Content,
		&respComment.CreatedAt,
		&updatedAT,
	); err != nil {
		return nil, err
	}

	if updatedAT.Valid {
		respComment.UpdatedAt = updatedAT.Time.String()
	}

	return &respComment, nil
}

func (c *CommentRepo) UpdateComment(req *pbc.Comment) (*pbc.Comment, error) {
	var updatedAT sql.NullTime
	query := `UPDATE comments 
			  SET content = $1
			  WHERE id = $2
			  RETURNING id, 
			  			user_id, 
			      		post_id,
			  			content, 
			  			created_at, 
			  			updated_at`
	var respPost pbc.Comment
	rowPost := c.db.QueryRow(query, req.Content, req.Id)
	if err := rowPost.Scan(&respPost.Id,
		&respPost.UserId,
		&respPost.PostId,
		&respPost.Content,
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

func (c *CommentRepo) DeleteComment(req *pbc.CommentId) (*pbc.Comment, error) {
	query := `DELETE FROM comments WHERE id = $1 RETURNING id, 
			  			user_id, 
    					post_id,
			  			content, 
			  			created_at, 
			  			updated_at`

	var respComment pbc.Comment
	var updatedAT sql.NullTime
	rowPost := c.db.QueryRow(query, req.Id)
	if err := rowPost.Scan(&respComment.Id,
		&respComment.UserId,
		&respComment.PostId,
		&respComment.Content,
		&respComment.CreatedAt,
		&updatedAT); err != nil {
		return nil, err
	}

	if updatedAT.Valid {
		respComment.UpdatedAt = updatedAT.Time.String()
	}

	return &respComment, nil
}

func (c *CommentRepo) GetAllComments(req *pbc.CommentsRequest) (*pbc.CommentsResponse, error) {
	query := `SELECT id, 
       				 user_id, 
       				 post_id,
       				 content
			  FROM comments 
			  WHERE post_id = $1`
	rows, err := c.db.Query(query, req.PostId)
	if err != nil {
		return nil, err
	}
	var response pbc.CommentsResponse

	for rows.Next() {
		var respComment pbc.Comment
		err := rows.Scan(&respComment.Id, &respComment.UserId, &respComment.PostId, &respComment.Content)
		if err != nil {
			return nil, err
		}

		response.Comments = append(response.Comments, &respComment)
	}

	return &response, nil
}
