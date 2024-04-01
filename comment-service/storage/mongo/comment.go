package mongodb

import (
	pbc "4microservices/comment-service/genproto/comment_service"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentRepo struct {
	collection *mongo.Collection
}

func NewCommentRepo(client *mongo.Client, dbName, collectionName string) *CommentRepo {
	collection := client.Database(dbName).Collection(collectionName)
	return &CommentRepo{
		collection: collection,
	}
}

func (c *CommentRepo) CreateComment(req *pbc.Comment) (*pbc.Comment, error) {
	_, err := c.collection.InsertOne(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *CommentRepo) GetComment(req *pbc.CommentId) (*pbc.Comment, error) {
	var comment pbc.Comment
	filter := bson.M{"id": req.Id}
	err := c.collection.FindOne(context.Background(), filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *CommentRepo) UpdateComment(req *pbc.Comment) (*pbc.Comment, error) {
	filter := bson.M{"id": req.Id}
	update := bson.M{"$set": bson.M{"content": req.Content}}
	_, err := c.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (c *CommentRepo) DeleteComment(req *pbc.CommentId) (*pbc.Comment, error) {
	var comment pbc.Comment
	filter := bson.M{"id": req.Id}
	err := c.collection.FindOneAndDelete(context.Background(), filter).Decode(&comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (c *CommentRepo) GetAllComments(req *pbc.CommentsRequest) (*pbc.CommentsResponse, error) {
	var comments []*pbc.Comment
	filter := bson.M{"post_id": req.PostId}
	cursor, err := c.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var comment pbc.Comment
		err := cursor.Decode(&comment)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &pbc.CommentsResponse{Comments: comments}, nil
}
