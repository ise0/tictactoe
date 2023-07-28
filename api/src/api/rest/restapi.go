package restapi

import (
	userRouter "api/src/api/rest/routes/user"
	wsapi "api/src/api/ws"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
)

var Engine = gin.Default()

func init() {
	Engine.Use(func(ctx *gin.Context) {
		cors.New(cors.Options{AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), " "),
			AllowCredentials: true}).HandlerFunc(ctx.Writer, ctx.Request)
	})
	r := Engine.Group("api")
	r.GET("/ws", wsapi.Upgrader)
	r.GET("/time", func(ctx *gin.Context) { ctx.JSON(200, time.Now().UnixMilli()) })

	userRouter.Apply(r.Group("/user"))
}
