//crud

package db

import (
	"context"
	"time"
)

type User struct {
	Username         string    `json:"username"`
	HashedPassword   string    `json:"hashed_password"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	PasswordChangeAt time.Time `json:"password_changed_at"`
	CreateAt         time.Time `json:"crate_at"`
}

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

const createUser = `-- name: CreateUser :one
INSERT INTO users(
	username,
	hashed_password,
	full_name,
	email,
)	VALUES (
	$1,$2,$3,$4
)

RETURNING username,hashed_password,full_name,email,password_changed_at,create_at;
`

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.HashedPassword, arg.FullName, arg.Email)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangeAt,
		&i.CreateAt,
	)

	return i, err
}

const getUser = `-- name :GetUser :one
SELECT username,hashed_password,full_name,email,password_changed_at,create_at FROM users
WHERE username =$1 LIMIT 1;
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.PasswordChangeAt,
		&i.CreateAt,
	)
	return i, err
}
