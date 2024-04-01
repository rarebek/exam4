package cassandra

import (
	pbc "4microservices/comment-service/genproto/comment_service"
	"github.com/gocql/gocql"
)

type CommentRepo struct {
	session *gocql.Session
}

func NewCommentRepo(clusterHosts []string, keyspace string) (*CommentRepo, error) {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &CommentRepo{session: session}, nil
}

func (c *CommentRepo) Close() {
	c.session.Close()
}

func (c *CommentRepo) CreateComment(req *pbc.Comment) (*pbc.Comment, error) {
	query := `INSERT INTO comments (id, user_id, post_id, content) VALUES (?, ?, ?, ?) USING TIMESTAMP now()`
	req.Id = gocql.TimeUUID().String()

	if err := c.session.Query(query, req.Id, req.UserId, req.PostId, req.Content).Exec(); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *CommentRepo) GetComment(req *pbc.CommentId) (*pbc.Comment, error) {
	query := `SELECT id, user_id, post_id, content FROM comments WHERE id = ?`

	var respComment pbc.Comment
	err := c.session.Query(query, req.Id).Scan(&respComment.Id, &respComment.UserId, &respComment.PostId, &respComment.Content)
	if err != nil {
		return nil, err
	}
	return &respComment, nil
}

func (c *CommentRepo) UpdateComment(req *pbc.Comment) (*pbc.Comment, error) {
	query := `UPDATE comments SET content = ? WHERE id = ? USING TIMESTAMP now()`
	if err := c.session.Query(query, req.Content, req.Id).Exec(); err != nil {
		return nil, err
	}
	return req, nil
}

func (c *CommentRepo) DeleteComment(req *pbc.CommentId) error {
	query := `DELETE FROM comments WHERE id = ? USING TIMESTAMP now()`
	if err := c.session.Query(query, req.Id).Exec(); err != nil {
		return err
	}
	return nil
}

func (c *CommentRepo) GetAllComments(req *pbc.CommentsRequest) (*pbc.CommentsResponse, error) {
	query := `SELECT id, user_id, post_id, content FROM comments WHERE post_id = ?`

	iter := c.session.Query(query, req.PostId).Iter()
	var response pbc.CommentsResponse
	for {
		var respComment pbc.Comment
		if !iter.Scan(&respComment.Id, &respComment.UserId, &respComment.PostId, &respComment.Content) {
			break
		}
		response.Comments = append(response.Comments, &respComment)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return &response, nil
}
