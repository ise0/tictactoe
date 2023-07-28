package userModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func GetById(id int) (user, error) {
	var u user
	if err := db.DB.QueryRow(context.TODO(), `
		select user_id, is_anonym, username, password from users where user_id = $1
	`, id).Scan(&u.Id, &u.IsAnonym, &u.Username, &u.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user{}, fmt.Errorf("user with id:%d does not exists", id)
		}
		logger.L.Errorw("Failed to execute sql query")
		return user{}, err
	}

	return u, nil
}
