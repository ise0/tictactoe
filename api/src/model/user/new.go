package userModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"api/src/shared/passhash"
	"context"
	"errors"
)

var ErrUsernameAlreadyExists = errors.New("username already exists")

func New(username string, password string) (user, error) {
	if len(password) < 1 {
		return user{}, errors.New("weak password")
	}
	if len(username) < 1 {
		return user{}, errors.New("short username")
	}
	var u user
	passwordHash, err := passhash.Hash(password)
	if err != nil {
		return user{}, err
	}

	if err := db.DB.QueryRow(context.TODO(), `
		insert into users(is_anonym, username, user_password) values (false, $1, $2)
		returning user_id, is_anonym, username, user_password
	`, username, passwordHash).Scan(&u.Id, &u.IsAnonym, &u.Username, &u.Password); err != nil {
		logger.L.Errorw("Failed to execute sql query")
		return user{}, err
	}

	return u, nil
}
