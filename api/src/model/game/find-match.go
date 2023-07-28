package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
)

func FindMatch(userId int) error {
	_, err := db.DB.Exec(context.TODO(), `
		insert into game_match_search values ($1)
		on conflict do nothing
	`, userId)
	if err != nil {
		logger.L.Errorw("Failed to execute sql query", "error", err)
		return err
	}

	return nil
}
