package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"

	"github.com/jackc/pgx/v5"
)

func CancelMatchSearch(users []int) error {
	sqlFailLog := func(err error) error {
		logger.L.Errorw("Failed to execute sql query", "error", err)
		return err
	}

	tx, err := db.DB.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer tx.Rollback(context.TODO())

	if err != nil {
		return sqlFailLog(err)
	}

	_, err = tx.Exec(context.TODO(), `delete from game_match_search where user_id = any($1)`, users)
	if err != nil {
		return sqlFailLog(err)
	}

	_, err = tx.Exec(context.TODO(), `delete from users where user_id = any($1) and is_anonym = true`, users)
	if err != nil {
		return sqlFailLog(err)
	}

	if err = tx.Commit(context.TODO()); err != nil {
		return sqlFailLog(err)
	}
	return nil
}
