package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub, userID string, roundID string, roomID string, phase string, letter string, timer int, gameTime int, w http.ResponseWriter, r *http.Request) {
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

	if phase != "" {
		type GameState struct {
			Phase    string `json:"phase"`
			Letter   string `json:"letter"`
			Timer    int    `json:"timer"`
			GameTime int    `json:"gameTime"`
			RoundID  string `json:"roundID"`
		}

		state := GameState{
			Phase:    phase,
			Letter:   letter,
			Timer:    timer,
			GameTime: gameTime,
			RoundID:  roundID,
		}

		stateJSON, err := json.Marshal(state)
		if err != nil {
			log.Printf("failed to marshal game state for user %s: %v", userID, err)
			return
		}
		client.Channel <- append([]byte("STATE:"), stateJSON...)
	}

	go client.WritePump()
	go client.ReadPump(hub)
}
