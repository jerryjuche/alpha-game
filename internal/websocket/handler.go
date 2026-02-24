package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub, userID string, roomID string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	client := &Client{
		webConn: conn,
		UserId:  userID,
		RoomId:  roomID,
		Channel: make(chan []byte, 256),
	}

	hub.NewClient <- client

	go client.WritePump()
	go client.ReadPump(hub)
}
