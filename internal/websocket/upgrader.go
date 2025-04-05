package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type Upgrader struct {
	Hub *Hub
}

func NewUpgrader(hub *Hub) *Upgrader {
	return &Upgrader{Hub: hub}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (u *Upgrader) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	client := &Client{
		hub:  u.Hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	u.Hub.register <- client

	go client.writePump()
	go client.readPump()
}
