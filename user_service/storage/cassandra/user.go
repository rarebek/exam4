package cassandra

import (
	pb "4microservice/user-service/genproto/user_service"
	"fmt"
	"github.com/gocql/gocql"
)

type UserRepo struct {
	session *gocql.Session
}

func NewUserRepo(clusterHosts []string, keyspace string) (*UserRepo, error) {
	cluster := gocql.NewCluster(clusterHosts...)
	cluster.Keyspace = keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	return &UserRepo{session: session}, nil
}

func (u *UserRepo) Close() {
	u.session.Close()
}

func (u *UserRepo) CreateUser(user *pb.User) (*pb.User, error) {
	query := `INSERT INTO users (id, username, email, password, first_name, last_name, bio, website, refresh_token) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	user.Id = gocql.TimeUUID().String()

	if err := u.session.Query(query, user.Id, user.Username, user.Email, user.Password, user.FirstName, user.LastName, user.Bio, user.Website, user.RefreshToken).Exec(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) GetUser(req *pb.UserIdd) (*pb.GetUserResponse, error) {
	var user pb.GetUserResponse
	query := `SELECT id, username, email, password, first_name, last_name, bio, website FROM users WHERE id = ?`
	err := u.session.Query(query, req.Id).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Bio, &user.Website)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) UpdateUser(user *pb.User) (*pb.User, error) {
	query := `UPDATE users SET username = ?, first_name = ?, last_name = ?, bio = ?, website = ? WHERE id = ?`
	if err := u.session.Query(query, user.Username, user.FirstName, user.LastName, user.Bio, user.Website, user.Id).Exec(); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepo) DeleteUser(id *pb.UserIdd) error {
	query := `DELETE FROM users WHERE id = ?`
	if err := u.session.Query(query, id.Id).Exec(); err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetAllUsers(req *pb.GetAllRequest) (*pb.GetAllResponse, error) {
	var response pb.GetAllResponse
	query := `SELECT id, username, email, password, first_name, last_name, bio, website, refresh_token FROM users LIMIT ? OFFSET ?`
	iter := u.session.Query(query, req.Limit, (req.Page-1)*req.Limit).Iter()
	var user pb.GetUserResponse
	for iter.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Bio, &user.Website, &user.RefreshToken) {
		response.Users = append(response.Users, &user)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return &response, nil
}

func (u *UserRepo) CheckField(req *pb.CheckFieldRequest) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE %s = ?", req.Field)

	var count int
	err := u.session.Query(query, req.Value).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (u *UserRepo) GetUserByEmail(req *pb.Email) (*pb.User, error) {
	var user pb.User
	query := `SELECT id, username, email, password, first_name, last_name, bio, website FROM users WHERE email = ?`
	err := u.session.Query(query, req.Email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.Bio, &user.Website)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
