package websocket

type BroadcastMessage struct {
	RoomId  string
	Message []byte
}

type Hub struct {
	ClientRoom   map[string][]*Client
	NewClient    chan *Client
	UnregClients chan *Client
	BroadcastMsg chan BroadcastMessage
}

func NewHub() *Hub {
	return &Hub{
		ClientRoom:   make(map[string][]*Client),
		NewClient:    make(chan *Client),
		UnregClients: make(chan *Client),
		BroadcastMsg: make(chan BroadcastMessage),
	}

}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.NewClient:
			h.ClientRoom[client.RoomId] = append(h.ClientRoom[client.RoomId], client)

		case client := <-h.UnregClients:
			remaining := []*Client{}
			for _, c := range h.ClientRoom[client.RoomId] {
				if c.UserId != client.UserId {
					remaining = append(remaining, c)
				}
			}
			h.ClientRoom[client.RoomId] = remaining

		case msg := <-h.BroadcastMsg:
			for _, client := range h.ClientRoom[msg.RoomId] {
				client.Channel <- msg.Message
			}
		}
	}
}
