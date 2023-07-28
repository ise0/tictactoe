package userRouter

import (
	userModel "api/src/model/user"
	"api/src/shared/authjwt"
	"api/src/shared/passhash"

	"github.com/gin-gonic/gin"
)

func Apply(g *gin.RouterGroup) {
	g.POST("/sign-in", login)
	g.POST("/sign-up", signUp)
}

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type response struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	IsAnonym bool   `json:"isAnonym"`
	Rating   int    `json:"rating"`
	Jwt      string `json:"jwt"`
}

func login(c *gin.Context) {
	var u user

	if err := c.BindJSON(&u); err != nil {
		return
	}

	user, err := userModel.GetByUsername(u.Username)
	if err != nil || !passhash.Verify(u.Password, user.Password) {
		c.JSON(500, `"wrong password or username"`)
		return
	}

	token, err := authjwt.SignAuthJwt(user.Id)
	if err != nil {
		c.JSON(500, `"something went wrong"`)
		return
	}

	res := response{user.Id, user.Username, user.IsAnonym, user.Rating, token}

	c.JSON(200, res)
}

func signUp(c *gin.Context) {
	var u user

	if err := c.BindJSON(&u); err != nil {
		return
	}

	user, err := userModel.New(u.Username, u.Password)
	if err != nil {
		c.JSON(500, `"oops"`)
		return
	}

	token, err := authjwt.SignAuthJwt(user.Id)
	if err != nil {
		c.JSON(500, `"oops"`)
		return
	}

	res := response{user.Id, user.Username, user.IsAnonym, user.Rating, token}

	c.JSON(200, res)
}
