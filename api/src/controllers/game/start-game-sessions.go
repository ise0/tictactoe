package gameController

import (
	"api/src/api/ws/conns"
	gameModel "api/src/model/game"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func StartGameSessions(newGameChan chan<- gameModel.Game, timeExpiredChan chan<- TimeExpiredWinner) {
	s, err := gameModel.GetNewSessions()
	if err != nil {
		return
	}

	for _, v := range s {
		go func(players []int) {
			var (
				offlinePlayers, onlinePlayers []int
				m                             sync.Mutex
				wg                            sync.WaitGroup
			)
			wg.Add(len(players))
			for _, v2 := range players {
				go func(player int) {
					err := conns.Players[player].WriteMessage(websocket.PingMessage, []byte{})
					m.Lock()
					if err != nil {
						offlinePlayers = append(onlinePlayers, player)
						conns.PlayersMutex.Lock()
						delete(conns.Players, player)
						conns.PlayersMutex.Unlock()
					} else {
						onlinePlayers = append(offlinePlayers, player)

					}
					m.Unlock()
					wg.Done()
				}(v2)
			}
			wg.Wait()

			if len(offlinePlayers) > 0 {
				gameModel.CancelMatchSearch(offlinePlayers)
				gameModel.BackToMatchSearch(onlinePlayers)
			} else {
				g, err := gameModel.New(players[0], players[1])
				if err == nil {
					newGameChan <- g
					t := time.AfterFunc(gameModel.GameTimer,
						func() {
							winner, err := gameModel.TimeExpired(g.GameSessionId)
							if err != nil {
								return
							}
							timeExpiredChan <- TimeExpiredWinner{winner, g.PlayerA.UserId, g.PlayerB.UserId}
						},
					)
					gameTimersMutex.Lock()
					gameTimers[g.GameSessionId] = t
					gameTimersMutex.Unlock()
				}
			}
		}(v)
	}
}
