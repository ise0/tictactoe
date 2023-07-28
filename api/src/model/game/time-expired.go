package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var GameTimer = time.Second*30 + time.Millisecond*500

func TimeExpired(gameSessionId int) (int, error) {
	sqlFailLog := func(err error) (int, error) {
		logger.L.Errorw("Failed to execute sql query", "error", err)
		return 0, err
	}

	tx, err := db.DB.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	defer tx.Rollback(context.TODO())

	if err != nil {
		return sqlFailLog(err)
	}

	var winner int
	var players []int
	err = tx.QueryRow(context.TODO(), `
		with g as (
			select current_turn_player_id, game_session_id 
			from game_sessions 
			where game_session_id = $1 and now() > last_move_ts + $2
		), players as (
			select user_id 
			from game_players 
			where game_session_id = (select game_session_id from g)
		)
		select user_id, (select array_agg(user_id) from players) 
		from players 
		where user_id != (select current_turn_player_id from g) 
	`, gameSessionId, fmt.Sprint(GameTimer/1000)).Scan(&winner, &players)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := errors.New("game has already changed")
			logger.L.Info(err.Error())
			return 0, err
		}
		return sqlFailLog(err)
	}

	_, err = tx.Exec(context.TODO(), ` 
		update users set rating = rating + 
			case when user_id = $1 then 100 
			when rating < 100 then -rating
			else -100
			end
		where is_anonym = false and user_id in (
			select user_id 
			from game_players 
			where game_session_id = $2
		) 
	`, winner, gameSessionId)
	if err != nil {
		return sqlFailLog(err)
	}

	_, err = tx.Exec(context.TODO(), ` 
		delete from game_sessions_ids where game_session_id = $1
	`, gameSessionId)
	if err != nil {
		return sqlFailLog(err)
	}

	_, err = tx.Exec(context.TODO(), ` 
		delete from users where user_id = any($1) and is_anonym = true
	`, players)
	if err != nil {
		return sqlFailLog(err)
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		return sqlFailLog(err)
	}
	return winner, nil
}
