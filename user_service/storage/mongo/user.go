package mongodb

import (
	pb "4microservice/user-service/genproto/user_service"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserMongoRepo struct {
	collection *mongo.Collection
}

// NewUserMongoRepo creates a new instance of UserRepo
func NewUserMongoRepo(client *mongo.Client, dbName, collectionName string) *UserMongoRepo {
	collection := client.Database(dbName).Collection(collectionName)
	return &UserMongoRepo{collection: collection}
}

func (u *UserMongoRepo) CreateUser(user *pb.User) (*pb.User, error) {
	_, err := u.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserMongoRepo) GetUser(id *pb.UserIdd) (*pb.GetUserResponse, error) {
	var userWithPosts pb.GetUserResponse

	// Lookup stage to join users collection with posts collection
	lookupStage := bson.D{{"$lookup", bson.D{
		{"from", "posts"},
		{"localField", "id"},
		{"foreignField", "userId"},
		{"as", "posts"},
	}}}

	// Match stage to filter based on user ID
	matchStage := bson.D{{"$match", bson.M{"id": id.Id}}}

	// Aggregation pipeline
	pipeline := mongo.Pipeline{lookupStage, matchStage}

	// Aggregate operation
	cursor, err := u.collection.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate over the cursor
	for cursor.Next(context.Background()) {
		err := cursor.Decode(&userWithPosts)
		if err != nil {
			return nil, err
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &userWithPosts, nil
}

func (u *UserMongoRepo) UpdateUser(user *pb.User) (*pb.User, error) {
	filter := bson.M{"id": user.Id}
	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"bio":        user.Bio,
			"website":    user.Website,
		},
	}
	_, err := u.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserMongoRepo) DeleteUser(id *pb.UserIdd) (*empty.Empty, error) {
	filter := bson.M{"id": id}
	_, _ = u.collection.DeleteOne(context.Background(), filter)
	return nil, nil
}

func (u *UserMongoRepo) GetAllUsers(req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	var response pb.GetAllResponse
	offset := (req.Page - 1) * req.Limit
	cursor, err := u.collection.Find(context.Background(), bson.D{}, options.Find().SetLimit(int64(req.Limit)).SetSkip(offset))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var user pb.GetUserResponse
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		response.Users = append(response.Users, &user)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &response, nil
}

func (u *UserMongoRepo) CheckField(req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	filter := bson.M{req.Field: req.Value}
	count, err := u.collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return &pb.CheckFieldResponse{Unique: false}, err
	}
	return &pb.CheckFieldResponse{Unique: count > 0}, nil
}

func (u *UserMongoRepo) GetUserByEmail(req *pb.Email) (*pb.User, error) {
	var user pb.User
	filter := bson.M{"email": req.Email}
	err := u.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
