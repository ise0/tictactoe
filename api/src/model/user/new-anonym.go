package userModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
)

func NewAnonym() (int, error) {
	var userId int
	if err := db.DB.QueryRow(context.TODO(), `
		insert into users(is_anonym) values (true)
		returning user_id
	`).Scan(&userId); err != nil {
		logger.L.Errorw("Failed to execute sql query")
		return userId, err
	}

	return userId, nil
}
