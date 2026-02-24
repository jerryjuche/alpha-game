package websocket

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	webConn *websocket.Conn
	UserId  string
	RoomId  string
	Channel chan []byte
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.UnregClients <- c
		c.webConn.Close()
	}()

	for {
		_, message, err := c.webConn.ReadMessage()
		if err != nil {
			break
		}
		hub.BroadcastMsg <- BroadcastMessage{
			RoomId:  c.RoomId,
			Message: message,
		}
	}
}

func (c *Client) WritePump() {
	defer c.webConn.Close()

	for message := range c.Channel {
		err := c.webConn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
