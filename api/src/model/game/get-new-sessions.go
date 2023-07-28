package gameModel

import (
	"api/src/shared/db"
	"api/src/shared/logger"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

func GetNewSessions() ([][]int, error) {
	sqlFailLog := func(err error) ([][]int, error) {
		logger.L.Errorw("Failed to execute sql query")
		return nil, err
	}
	tx, err := db.DB.BeginTx(context.TODO(), pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	defer tx.Rollback(context.TODO())

	if err != nil {
		return sqlFailLog(err)
	}
	var res []byte
	if err := tx.QueryRow(context.TODO(), `
		with u as (
			update game_match_search set pending = true 
			where user_id in (
				select user_id from game_match_search 
				where pending = false
				offset case when (select count(*) from game_match_search ) % 2 = 0 then 0 else 1 end
			)
			returning user_id
		), p as (
			select user_id, ntile(round((select count(*)/2::numeric from u))::int) over(order by user_id) as session_n
			from u
		) 
		select to_json(array(select array_agg(user_id) 
		from p
		group by session_n))
	`).Scan(&res); err != nil {
		return sqlFailLog(err)
	}

	var wp [][]int
	if err = json.Unmarshal(res, &wp); err != nil {
		logger.L.Errorw("Failed to parse json", "error", err)
		return nil, err
	}

	if err = tx.Commit(context.TODO()); err != nil {
		return sqlFailLog(err)
	}
	return wp, nil
}
