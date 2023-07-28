package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type message struct {
	ChatMembers []int
	UserId      int
	Text        string
	Timestamp   time.Time
}

func AddMessage(userId int, text string) (message, error) {
	var res []byte
	ts := time.Now()
	if err := db.DB.QueryRow(context.TODO(), `
		with game_id as (
			select game_session_id 
			from game_players 
			where user_id = $1
		), game_id2 as (
			insert into game_chat(game_session_id, user_id, ts, msg) values
			((select game_session_id from game_id), $1, $2, $3)
			returning game_session_id
		) 
		select to_json(array_agg(user_id)) from game_players 
		where game_session_id = (select game_session_id from game_id2)
	`, userId, ts, text).Scan(&res); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := errors.New("game is over or user does not have the required rights")
			logger.L.Info(err.Error())
			return message{}, err
		}
		logger.L.Errorw("Failed to execute sql query", "error", err)
		return message{}, err
	}
	var cm []int
	if err := json.Unmarshal(res, &cm); err != nil {
		logger.L.Errorw("Failed to parse json", "error", err)
		return message{}, err
	}

	return message{cm, userId, text, ts}, nil
}
