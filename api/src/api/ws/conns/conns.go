package conns

import (
	"sync"

	"github.com/gorilla/websocket"
)

var Players = make(map[int]*websocket.Conn)
var PlayersMutex sync.Mutex
