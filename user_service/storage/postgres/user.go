package postgres

import (
	pb "4microservice/user-service/genproto/user_service"
	"database/sql"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type UserRepo struct {
	db   *sql.DB
	conn *nats.Conn
}

// NewUserRepo ...
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (u *UserRepo) CreateUser(user *pb.User) (*pb.User, error) {
	query := `INSERT INTO users(id, username, email, password, first_name, last_name, bio, website, refresh_token) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, username, email, password, first_name, last_name, bio, website, refresh_token`
	user.Id = uuid.New().String()

	var respUser pb.User

	rowUser := u.db.QueryRow(query, user.Id, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.Bio, user.Website, user.RefreshToken)
	if err := rowUser.Scan(&respUser.Id, &respUser.Username, &respUser.Email, &respUser.Password, &respUser.FirstName, &respUser.LastName, &respUser.Bio, &respUser.Website, &respUser.RefreshToken); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (u *UserRepo) GetUser(req *pb.UserIdd) (*pb.GetUserResponse, error) {
	query := `SELECT id, username, email, password, first_name, last_name, bio, website FROM users WHERE id = $1`

	var respUser pb.GetUserResponse

	rowUser := u.db.QueryRow(query, req.Id)
	if err := rowUser.Scan(&respUser.Id, &respUser.Username, &respUser.Email, &respUser.Password, &respUser.FirstName, &respUser.LastName, &respUser.Bio, &respUser.Website); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (u *UserRepo) UpdateUser(user *pb.User) (*pb.User, error) {
	query := `UPDATE users SET username = $1, first_name = $2, last_name = $3, bio = $4, website = $5 WHERE id = $6  RETURNING id, username, email, password, first_name, last_name, bio, website, refresh_token`

	var respUser pb.User

	rowUser := u.db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Bio, user.Website, user.Id)
	if err := rowUser.Scan(&respUser.Id, &respUser.Username, &respUser.Email, &respUser.Password, &respUser.FirstName, &respUser.LastName, &respUser.Bio, &respUser.Website, &respUser.RefreshToken); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func (u *UserRepo) DeleteUser(id *pb.UserIdd) (*empty.Empty, error) {
	query := `DELETE FROM users WHERE id = $1`
	u.db.QueryRow(query, id.Id)
	return &empty.Empty{}, nil
}

func (u *UserRepo) GetAllUsers(req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	var response pb.GetAllResponse
	offset := (req.Page - 1) * req.Limit
	query := `SELECT id, username, email, password, first_name, last_name, bio, website, refresh_token FROM users LIMIT $1 OFFSET $2`
	rows, err := u.db.Query(query, req.Limit, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			user = pb.GetUserResponse{Posts: []*pb.Postt{}}
		)
		err := rows.Scan(&user.Id, &user.Email, &user.Password, &user.Username, &user.Bio, &user.FirstName, &user.LastName, &user.Website, &user.RefreshToken)
		if err != nil {
			return nil, err
		}

		response.Users = append(response.Users, &user)
	}

	return &response, nil
}

func (u *UserRepo) CheckField(req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM users WHERE %s = $1)", req.Field)

	var exists bool
	err := u.db.QueryRow(query, req.Value).Scan(&exists)
	if err != nil {
		return nil, err
	}

	return &pb.CheckFieldResponse{Unique: exists}, nil
}

func (u *UserRepo) GetUserByEmail(req *pb.Email) (*pb.User, error) {
	query := `SELECT id, username, email, password, first_name, last_name, bio, website FROM users WHERE email = $1`

	var respUser pb.User

	rowUser := u.db.QueryRow(query, req.Email)
	if err := rowUser.Scan(&respUser.Id, &respUser.Username, &respUser.Email, &respUser.Password, &respUser.FirstName, &respUser.LastName, &respUser.Bio, &respUser.Website); err != nil {
		return nil, err
	}

	return &respUser, nil
}
