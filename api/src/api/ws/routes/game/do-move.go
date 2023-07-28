package gameRouter

import (
	"api/src/api/ws/conns"
	gameController "api/src/controllers/game"
	gameModel "api/src/model/game"
	"api/src/shared/authjwt"
	"api/src/shared/logger"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type doMoveAction struct {
	Value struct {
		Jwt    string `json:"jwt"`
		Coords struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"coords"`
	} `json:"value"`
}

func doMove(c *websocket.Conn, m []byte) {
	var a doMoveAction
	if err := json.Unmarshal(m, &a); err != nil {
		logger.L.Errorw("Failed to parse json ws", "error", err)
		return
	}

	userId, err := authjwt.ParseAuthJwt(a.Value.Jwt)
	if err != nil {
		return
	}

	b, err := gameController.DoMove(userId, gameModel.Coords(a.Value.Coords))
	if err != nil {
		return
	}

	newGameBoardMsg := NewGameBoardMsg(GameBoardData{b.Winner, b.Board, b.CurrentTurnPlayerId, b.LastMoveTs.UnixMilli()})

	if conns.Players[b.PlayerA] != nil {
		go conns.Players[b.PlayerA].WriteJSON(newGameBoardMsg)
	}
	if conns.Players[b.PlayerB] != nil {
		go conns.Players[b.PlayerB].WriteJSON(newGameBoardMsg)
	}
	if b.Winner != 0 {
		delete(conns.Players, b.PlayerA)
		delete(conns.Players, b.PlayerB)
	}
}
