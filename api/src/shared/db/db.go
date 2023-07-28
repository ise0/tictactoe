package db

import (
	"api/src/shared/logger"
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *pgx.Conn

func Connect() error {
	d, err := pgx.Connect(context.TODO(), os.Getenv("DB_CONNECTION"))

	if err != nil {
		logger.L.Errorw("Failed to connect to db", "error", err)
		return err
	}

	if _, err = d.Exec(context.TODO(), `delete from game_sessions_ids`); err != nil {
		logger.L.Errorw("Failed to ensure data consistency ", "error", err)
		return err
	}

	if _, err = d.Exec(context.TODO(), `delete from game_match_search`); err != nil {
		logger.L.Errorw("Failed to ensure data consistency ", "error", err)
		return err
	}

	if _, err = d.Exec(context.TODO(), `delete from users where is_anonym = true`); err != nil {
		logger.L.Errorw("Failed to ensure data consistency ", "error", err)
		return err
	}

	DB = d
	return nil
}
