package gameRouter

type (
	AnonymJwtData struct {
		UserId int    `json:"userId"`
		Jwt    string `json:"jwt"`
	}

	AnonymJwtMsg struct {
		Type string        `json:"type"`
		Data AnonymJwtData `json:"data"`
	}
)

func NewAnonymJwtMsg(data AnonymJwtData) AnonymJwtMsg {
	return AnonymJwtMsg{"game/anonJwt", data}
}

type (
	GameBoardData struct {
		Winner              int        `json:"winner"`
		Board               [][]string `json:"board"`
		CurrentTurnPlayerId int        `json:"currentTurnPlayerId"`
		LastMoveTs          int64      `json:"lastMoveTs"`
	}

	GameBoardMsg struct {
		Type string        `json:"type"`
		Data GameBoardData `json:"data"`
	}
)

func NewGameBoardMsg(data GameBoardData) GameBoardMsg {
	return GameBoardMsg{"game/gameBoard", data}
}

type (
	GameChatMsgData struct {
		UserId int    `json:"userId"`
		Ts     int64  `json:"ts"`
		Text   string `json:"text"`
	}
	GamePlayerData struct {
		UserId     int    `json:"userId"`
		IsAnonym   bool   `json:"isAnonym"`
		Username   string `json:"username"`
		PlayerSign string `json:"playerSign"`
		Rating     int    `json:"rating"`
	}
	GameData struct {
		GameSessionId       int               `json:"gameSessionId"`
		Board               [][]string        `json:"board"`
		Chat                []GameChatMsgData `json:"chat"`
		CurrentTurnPlayerId int               `json:"currentTurnPlayerId"`
		PlayerA             GamePlayerData    `json:"playerA"`
		PlayerB             GamePlayerData    `json:"playerB"`
		LastMoveTs          int64             `json:"lastMoveTs"`
	}
	GameMsg struct {
		T    string   `json:"type"`
		Data GameData `json:"data"`
	}
)

func NewGameMsg(data GameData) GameMsg {
	return GameMsg{"game/game", data}
}

type (
	ChatMsgData struct {
		UserId    int    `json:"userId"`
		Text      string `json:"text"`
		Timestamp int64  `json:"ts"`
	}

	ChatMsg struct {
		Type string      `json:"type"`
		Data ChatMsgData `json:"data"`
	}
)

func NewChatMsg(data ChatMsgData) ChatMsg {
	return ChatMsg{"game/chatMsg", data}
}

type MatchSearchStartedMsg struct {
	Type string `json:"type"`
}

func NewMatchSearchStartedMsg() MatchSearchStartedMsg {
	return MatchSearchStartedMsg{"game/matchSearchStarted"}
}

type (
	GameTimeExpiredData struct {
		Winner int `json:"winner"`
	}
	GameTimeExpiredMsg struct {
		Type string              `json:"type"`
		Data GameTimeExpiredData `json:"data"`
	}
)

func NewTimeExpiredMsg(winner int) GameTimeExpiredMsg {
	return GameTimeExpiredMsg{"game/timeExpired", GameTimeExpiredData{winner}}
}
