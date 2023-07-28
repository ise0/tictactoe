package gameModel

import "time"

type (
	chat_msg struct {
		UserId int       `json:"user_id"`
		Ts     time.Time `json:"ts"`
		Text   string    `json:"msg"`
	}
	player struct {
		UserId     int    `json:"user_id"`
		IsAnonym   bool   `json:"is_anonym"`
		Username   string `json:"username"`
		PlayerSign string `json:"player_sign"`
		Rating     int    `json:"rating"`
	}
	Game struct {
		GameSessionId       int        `json:"game_session_id"`
		Board               [][]string `json:"board"`
		CurrentTurnPlayerId int        `json:"current_turn_player_id"`
		LastMoveTs          time.Time  `json:"last_move_ts"`
		Chat                []chat_msg `json:"chat"`
		PlayerA             player     `json:"player_a"`
		PlayerB             player     `json:"player_b"`
	}
)
