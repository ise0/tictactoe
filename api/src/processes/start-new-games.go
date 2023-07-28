package proccesses

import (
	"api/src/api/ws/conns"
	gameRouter "api/src/api/ws/routes/game"
	gameController "api/src/controllers/game"
	gameModel "api/src/model/game"
	"time"
)

func createGame(g gameModel.Game) gameRouter.GameData {
	chat := make([]gameRouter.GameChatMsgData, 0, len(g.Chat))
	for _, c := range g.Chat {
		chat = append(chat, gameRouter.GameChatMsgData{UserId: c.UserId, Ts: c.Ts.UnixMilli(), Text: c.Text})
	}
	return gameRouter.GameData{
		GameSessionId:       g.GameSessionId,
		Board:               g.Board,
		Chat:                chat,
		CurrentTurnPlayerId: g.CurrentTurnPlayerId,
		PlayerA:             gameRouter.GamePlayerData(g.PlayerA),
		PlayerB:             gameRouter.GamePlayerData(g.PlayerB),
		LastMoveTs:          g.LastMoveTs.UnixMilli(),
	}
}

var newGameChan = make(chan gameModel.Game)
var gameTimeExpiredChan = make(chan gameController.TimeExpiredWinner)

func Start() {
	go StartNewGameSessions()
	go func() {
		for {
			g := <-gameTimeExpiredChan

			go func() {
				m := gameRouter.NewTimeExpiredMsg(g.Winner)
				go conns.Players[g.PlayerA].WriteJSON(m)
				go conns.Players[g.PlayerB].WriteJSON(m)
			}()
		}
	}()

	go func() {
		for {
			g := <-newGameChan

			go func() {
				m := gameRouter.NewGameMsg(createGame(g))
				go conns.Players[g.PlayerA.UserId].WriteJSON(m)
				go conns.Players[g.PlayerB.UserId].WriteJSON(m)
			}()
		}
	}()

}

func StartNewGameSessions() {
	for range time.Tick(time.Second * 5) {
		go gameController.StartGameSessions(newGameChan, gameTimeExpiredChan)
	}
}
