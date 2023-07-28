package gameRouter

import (
	"api/src/api/ws/conns"
	gameModel "api/src/model/game"
	"api/src/shared/authjwt"
	"api/src/shared/logger"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type getGameAction struct {
	Type string
	Data struct {
		Jwt string `json:"jwt"`
	} `json:"data"`
}

func get(c *websocket.Conn, m []byte) {
	var a getGameAction
	if err := json.Unmarshal(m, &a); err != nil {
		logger.L.Errorw("Failed to parse json ws", "error", err)
		return
	}
	userId, err := authjwt.ParseAuthJwt(a.Data.Jwt)
	if err != nil {
		return
	}
	g, err := gameModel.GetByUserId(userId)

	if err != nil {
		return
	}

	conns.Players[userId] = c

	go conns.Players[userId].WriteJSON(NewGameMsg(createGame(g)))
}

func createGame(g gameModel.Game) GameData {
	chat := make([]GameChatMsgData, 0, len(g.Chat))
	for _, c := range g.Chat {
		chat = append(chat, GameChatMsgData{c.UserId, c.Ts.UnixMilli(), c.Text})
	}
	return GameData{
		GameSessionId:       g.GameSessionId,
		Board:               g.Board,
		Chat:                chat,
		CurrentTurnPlayerId: g.CurrentTurnPlayerId,
		PlayerA:             GamePlayerData(g.PlayerA),
		PlayerB:             GamePlayerData(g.PlayerB),
		LastMoveTs:          g.LastMoveTs.UnixMilli(),
	}
}
