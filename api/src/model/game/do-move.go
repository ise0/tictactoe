package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Coords struct {
	X int
	Y int
}

type NewGameBoard struct {
	GameSessionId       int        `json:"game_session_id"`
	LastMoveTs          time.Time  `json:"last_move_ts"`
	Winner              int        `json:"winner"`
	Board               [][]string `json:"board"`
	CurrentTurnPlayerId int        `json:"current_turn_player_id"`
	PlayerA             int        `json:"player_a"`
	PlayerB             int        `json:"player_b"`
}

func DoMove(userId int, coords Coords) (NewGameBoard, error) {
	sqlFailLog := func(err error) (NewGameBoard, error) {
		logger.L.Errorw("Failed to execute sql query", "error", err)
		return NewGameBoard{}, err
	}

	tx, err := db.DB.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer tx.Rollback(context.TODO())

	if err != nil {
		return sqlFailLog(err)
	}

	var board NewGameBoard
	var gameBoardJSON []byte
	if err = tx.QueryRow(context.TODO(), `
		with players as (
			select * 
			from game_players 
			where game_session_id = (select game_session_id from game_players where user_id = $1)
		)
		update game_sessions set 
			game_board = jsonb_set(game_board, $4, (select to_jsonb(player_sign) from players where user_id = $1), false),
			current_turn_player_id = (select user_id from players where user_id != $1),
			last_move_ts = now()
		where game_session_id = (select game_session_id from players limit 1) 
			and current_turn_player_id = $1 and game_board->$3::int->>$2::int = ''
		returning game_session_id, current_turn_player_id, last_move_ts, game_board, (select user_id from players limit 1) as player_a, (select user_id from players offset 1) as player_b
	`, userId, coords.X, coords.Y, fmt.Sprintf("{%d, %d}", coords.Y, coords.X)).
		Scan(&board.GameSessionId, &board.CurrentTurnPlayerId, &board.LastMoveTs, &gameBoardJSON, &board.PlayerA, &board.PlayerB); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := errors.New("game is over or user does not have the required rights")
			logger.L.Info(err.Error())
			return NewGameBoard{}, err
		}
		return sqlFailLog(err)
	}
	if err = json.Unmarshal(gameBoardJSON, &board.Board); err != nil {
		logger.L.Errorw("Failed to parse json", "error", err)
		return NewGameBoard{}, err
	}

	if winnerSign := getWinner(board.Board); winnerSign != "" {
		board.Winner = userId

		if _, err := tx.Exec(context.TODO(), ` 
			update users set rating = rating + 
				case when user_id = $1 then 100 
				when rating < 100 then -rating
				else -100
				end
			where is_anonym = false and user_id in ($2, $3) 
		`, board.Winner, board.PlayerA, board.PlayerB); err != nil {
			return sqlFailLog(err)
		}

		if _, err := tx.Exec(context.TODO(), ` 
			delete from game_sessions_ids where game_session_id = $1
		`, board.GameSessionId); err != nil {
			return sqlFailLog(err)
		}

		if _, err := tx.Exec(context.TODO(), ` 
			delete from users 
			where is_anonym = true and user_id in ($1, $2) 
		`, board.PlayerA, board.PlayerB); err != nil {
			return sqlFailLog(err)
		}
	}

	if err = tx.Commit(context.TODO()); err != nil {
		return sqlFailLog(err)
	}
	return board, nil
}
