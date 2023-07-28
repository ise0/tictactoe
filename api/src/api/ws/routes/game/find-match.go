package gameRouter

import (
	"api/src/api/ws/conns"
	gameModel "api/src/model/game"
	userModel "api/src/model/user"
	"api/src/shared/authjwt"
	"api/src/shared/logger"
	"encoding/json"

	"github.com/gorilla/websocket"
)

type findMatchAction struct {
	Type string
	Data struct {
		Jwt      string `json:"jwt"`
		IsAnonym bool   `json:"isAnonym"`
	} `json:"data"`
}

func findMatch(c *websocket.Conn, m []byte) {
	var (
		a      findMatchAction
		userId int
		err    error
	)
	if err := json.Unmarshal(m, &a); err != nil {
		logger.L.Errorw("Failed to parse json ws", "error", err)
		return
	}
	if !a.Data.IsAnonym {
		if userId, err = authjwt.ParseAuthJwt(a.Data.Jwt); err != nil {
			return
		}
	} else {
		if userId, err = userModel.NewAnonym(); err != nil {
			return
		}
	}

	conns.PlayersMutex.Lock()
	conns.Players[userId] = c
	conns.PlayersMutex.Unlock()

	if err := gameModel.FindMatch(userId); err != nil {
		return
	}

	if a.Data.IsAnonym {
		t, err := authjwt.SignAuthJwt(userId)
		if err != nil {
			return
		}
		c.WriteJSON(NewAnonymJwtMsg(AnonymJwtData{userId, t}))
	}

	c.WriteJSON(NewMatchSearchStartedMsg())
}
