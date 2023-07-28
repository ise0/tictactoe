package gameRouter

import (
	"github.com/gorilla/websocket"
)

func Init(r *map[string]func(c *websocket.Conn, m []byte)) {
	(*r)["game/find-match"] = findMatch
	(*r)["game/cancel-match-search"] = cancelMatchSearch
	(*r)["game/do-move"] = doMove
	(*r)["game/get"] = get
	(*r)["game/add-chat-msg"] = addChatMessage
}
