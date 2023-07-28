package wsapi

import (
	gameRouter "api/src/api/ws/routes/game"
	"api/src/shared/logger"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var router = make(map[string]func(c *websocket.Conn, m []byte))

func init() {
	gameRouter.Init(&router)
}

var wsUpgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func Upgrader(ctx *gin.Context) {

	c, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		logger.L.Errorw("Failed to upgrade connection to ws", "error", err)
		return
	}
	for {
		err = read(c)
		if err != nil {
			break
		}
	}
}

type action struct {
	Type string `json:"type"`
}

func read(c *websocket.Conn) error {
	t, m, err := c.ReadMessage()

	if err != nil {
		logger.L.Errorw("Failed to read ws", "error", err)
		return err
	}
	if t != websocket.TextMessage {
		return nil
	}
	var a action
	err = json.Unmarshal(m, &a)
	if err != nil {
		logger.L.Errorw("Failed to parse json", "error", err)
		return err
	}

	if actionHandler, ok := router[a.Type]; ok {
		actionHandler(c, m)
	}

	return nil
}
