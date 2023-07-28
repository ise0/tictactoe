package gameRouter

import (
	"api/src/api/ws/conns"
	gameModel "api/src/model/game"
	"api/src/shared/authjwt"
	"api/src/shared/logger"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type addChatMsgAction struct {
	Value struct {
		Jwt  string `json:"jwt"`
		Text string `json:"text"`
	} `json:"data"`
}

func addChatMessage(c *websocket.Conn, m []byte) {
	var a addChatMsgAction
	if err := json.Unmarshal(m, &a); err != nil {
		logger.L.Errorw("Failed to parse json ws", "error", err)
		return
	}
	userId, err := authjwt.ParseAuthJwt(a.Value.Jwt)
	if err != nil {
		return
	}
	msg, err := gameModel.AddMessage(userId, a.Value.Text)
	if err != nil {
		return
	}

	conns.Players[userId] = c

	newChatMsg := NewChatMsg(ChatMsgData{msg.UserId, msg.Text, msg.Timestamp.Unix()})

	for _, v := range msg.ChatMembers {
		if conns.Players[v] != nil {
			go conns.Players[v].WriteJSON(newChatMsg)
		}
	}
}
