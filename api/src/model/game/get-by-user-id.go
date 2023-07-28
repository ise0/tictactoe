package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

func GetByUserId(userId int) (Game, error) {
	var res []byte

	err := db.DB.QueryRow(context.TODO(), `
		with game_id as (
			select game_session_id from game_players where user_id = $1
		)
		select to_json(get_game((select game_session_id from game_id)))
	`, userId).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) || len(res) == 0 {
		err := errors.New("game is over or user does not have the required rights")
		logger.L.Info(err.Error())
		return Game{}, err
	} else if err != nil {
		logger.L.Errorw("Failed to execute sql query", "error", err)
		return Game{}, err
	}
	var g Game
	if err = json.Unmarshal(res, &g); err != nil {
		logger.L.Errorw("Failed to parse json", "error", err)
		return Game{}, err
	}
	return g, nil
}
