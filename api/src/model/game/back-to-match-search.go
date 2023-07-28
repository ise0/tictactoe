package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
)

func BackToMatchSearch(users []int) error {

	_, err := db.DB.Exec(context.TODO(), `
		update game_match_search set pending = false where user_id = any($1)
		`, users)
	if err != nil {
		logger.L.Errorw("Failed to execute sql query", "error", err)
		return err
	}

	return nil
}
