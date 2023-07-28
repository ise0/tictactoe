package gameRouter

import (
	"api/src/api/ws/conns"
	gameModel "api/src/model/game"
	"api/src/shared/authjwt"
	"api/src/shared/logger"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type cancelMatchSearchAction struct {
	Type string
	Data struct {
		Jwt string `json:"jwt"`
	} `json:"data"`
}

func cancelMatchSearch(c *websocket.Conn, m []byte) {
	var a cancelMatchSearchAction
	if err := json.Unmarshal(m, &a); err != nil {
		logger.L.Errorw("Failed to parse json ws", "error", err)
		return
	}
	userId, err := authjwt.ParseAuthJwt(a.Data.Jwt)
	if err != nil {
		return
	}
	if err := gameModel.CancelMatchSearch([]int{userId}); err != nil {
		return
	}
	delete(conns.Players, userId)
}
