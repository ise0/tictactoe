package userModel

type user struct {
	Id       int    `json:"user_id"`
	IsAnonym bool   `json:"is_anonym"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	Password string `json:"password"`
}
