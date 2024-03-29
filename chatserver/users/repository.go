package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type User struct {
	Id       string `json:"id" pgx:"id"`
	Username string `json:"username" pgx:"username" binding:"required"`
	Email    string `json:"email" pgx:"email" binding:"required"`
}

type UserRepo struct {
	conn *pgx.Conn
}

func NewUserRepo(conn *pgx.Conn) UserRepo {
	return UserRepo{conn: conn}
}

func (r *UserRepo) CreateUser(user User) (User, error) {
	newId, err := uuid.NewV7()
	if err != nil {
		return User{}, errors.New("failed to generate new user id")
	}
	user.Id = newId.String()
	_, err = r.conn.Exec(context.Background(), "INSERT INTO app_user (id, username, email) VALUES ($1, $2, $3)", user.Id, user.Username, user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *UserRepo) GetUsers() ([]User, error) {
	rows, err := r.conn.Query(context.Background(), "SELECT id, username, email FROM app_user")
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, 3)
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.Id, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)

	}

	return users, nil
}

func (r *UserRepo) GetUser(id string) (User, bool) {
	row := r.conn.QueryRow(context.Background(), "SELECT id, username, email FROM app_user WHERE id=$1", id)
	user := User{}
	err := row.Scan(&user.Id, &user.Username, &user.Email)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return User{}, false
	}
	return user, true
}

func (r *UserRepo) UpdateUser(id string, user User) (User, error) {
	user.Id = id
	_, err := r.conn.Exec(context.Background(), "UPDATE app_user SET username=$2, email=$3 WHERE id=$1", id, user.Username, user.Email)
	if err != nil {
		return User{}, err
	}
	return user, nil

}

func (r *UserRepo) DeleteUser(id string) error {
	_, err := r.conn.Exec(context.Background(), "DELETE FROM app_user WHERE id=$1", id)
	return err
}
