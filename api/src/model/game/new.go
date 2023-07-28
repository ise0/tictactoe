package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

func New(userAId int, userBId int) (Game, error) {
	sqlFailLog := func(err error) (Game, error) {
		logger.L.Errorw("Failed to execute sql query")
		return Game{}, err
	}

	tx, err := db.DB.Begin(context.TODO())
	defer tx.Rollback(context.TODO())

	if err != nil {
		return sqlFailLog(err)
	}

	var res []byte
	if err = tx.QueryRow(context.TODO(), `
		with session_id as (
			insert into game_sessions_ids values (default) returning game_session_id
		), players as (
			insert into game_players (user_id, game_session_id, player_sign) values 
			($1, (select game_session_id from session_id), 'X'),
			($2, (select game_session_id from session_id), 'O')
			returning user_id
		), game_session as (
			insert into game_sessions values 
			((select game_session_id from session_id), default, (select user_id from players offset 1))
			returning game_session_id
		)
		select to_json(get_game((select game_session_id from game_session)))
	`, userAId, userBId).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Game{}, errors.New("game is over or user does not have the required rights")
		}
		return sqlFailLog(err)
	}

	var game Game
	if err = json.Unmarshal(res, &game); err != nil {
		logger.L.Errorw("Failed to parse json", "error", err)
		return Game{}, err
	}

	if _, err := tx.Exec(context.TODO(), `
		delete from game_match_search where user_id in ($1, $2)
	`, userAId, userBId); err != nil {
		return sqlFailLog(err)
	}

	if err = tx.Commit(context.TODO()); err != nil {
		return sqlFailLog(err)
	}

	return game, nil
}
