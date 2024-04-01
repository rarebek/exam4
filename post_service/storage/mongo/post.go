package mongodb

import (
	pb "4microservice/post-service/genproto/post_service"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepo struct {
	collection *mongo.Collection
}

func NewPostRepo(client *mongo.Client, dbName, collectionName string) *PostRepo {
	collection := client.Database(dbName).Collection(collectionName)
	return &PostRepo{collection: collection}
}

func (p *PostRepo) CreatePost(req *pb.Post) (*pb.Post, error) {
	req.Id = uuid.NewString()
	_, err := p.collection.InsertOne(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (p *PostRepo) GetPost(req *pb.PostID) (*pb.Post, error) {
	var post pb.Post
	filter := bson.M{"id": req.Id}
	err := p.collection.FindOne(context.Background(), filter).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p *PostRepo) UpdatePost(req *pb.Post) (*pb.Post, error) {
	filter := bson.M{"id": req.Id}
	update := bson.M{
		"$set": bson.M{
			"content":   req.Content,
			"image_url": req.ImageUrl,
			"title":     req.Title,
		},
	}
	_, err := p.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func (p *PostRepo) DeletePost(req *pb.PostID) (*empty.Empty, error) {
	var post pb.Post
	filter := bson.M{"id": req.Id}
	err := p.collection.FindOneAndDelete(context.Background(), filter).Decode(&post)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (p *PostRepo) GetAllPosts(req *pb.PostsRequest) (*pb.UserWithPosts, error) {
	var posts []*pb.Post
	options := options.Find().SetLimit(req.Limit).SetSkip((req.Page - 1) * req.Limit)
	cursor, err := p.collection.Find(context.Background(), bson.D{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var post pb.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &pb.UserWithPosts{Posts: posts}, nil
}

func (p *PostRepo) GetUserPostsByUserId(req *pb.UserId) (*pb.UserWithPosts, error) {
	var posts []*pb.Post
	filter := bson.M{"user_id": req.UserId}
	cursor, err := p.collection.Find(context.Background(), filter)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var post pb.Post
		err := cursor.Decode(&post)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &pb.UserWithPosts{Posts: posts}, nil
}
