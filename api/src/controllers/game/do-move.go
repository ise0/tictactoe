package gameController

import (
	gameModel "api/src/model/game"
)

func DoMove(userId int, coords gameModel.Coords) (gameModel.NewGameBoard, error) {
	g, err := gameModel.DoMove(userId, coords)
	if err != nil {
		return g, err
	}
	if g.Winner != 0 {
		gameTimers[g.GameSessionId].Stop()
		go func() {
			gameTimersMutex.Lock()
			delete(gameTimers, g.GameSessionId)
			gameTimersMutex.Unlock()
		}()
	} else {
		gameTimers[g.GameSessionId].Reset(gameModel.GameTimer)
	}

	return g, err
}
