package userModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func GetByUsername(username string) (user, error) {
	var u user
	if err := db.DB.QueryRow(context.TODO(), `
		select user_id, is_anonym, username, rating, user_password from users where username = $1
	`, username).Scan(&u.Id, &u.IsAnonym, &u.Username, &u.Rating, &u.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user{}, fmt.Errorf("user with username:%s does not exists", username)
		}
		logger.L.Errorw("Failed to execute sql query")
		return user{}, err
	}

	return u, nil
}
