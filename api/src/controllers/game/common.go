package gameController

import (
	"sync"
	"time"
)

type TimeExpiredWinner struct {
	Winner  int
	PlayerA int
	PlayerB int
}

var gameTimers = make(map[int]*time.Timer)
var gameTimersMutex sync.Mutex
