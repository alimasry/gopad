// this code is heavily inspired by the gorilla websocket chat example
// https://github.com/gorilla/websocket/blob/main/examples/chat/hub.go

package ws

import "sync"

type Hub struct {
	clients    ClientList
	broadcast  chan *Event
	register   chan *Client
	unregister chan *Client
}

var (
	hubInstance *Hub
	once        sync.Once
)

func InitializeHub() {
	once.Do(func() {
		hubInstance = &Hub{
			clients:    make(ClientList),
			broadcast:  make(chan *Event),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
	})
}

func GetHubInstance() *Hub {
	InitializeHub()
	return hubInstance
}

func (h *Hub) Run() {
	InitializeHub()
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case event := <-h.broadcast:
			for client := range h.clients[event.UUID] {
				select {
				case client.send <- event:
				default:
					h.removeClient(client)
				}
			}
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client.UUID][client]; ok {
		delete(h.clients[client.UUID], client)
		close(client.send)

		if len(h.clients) == 0 {
			delete(h.clients, client.UUID)
		}
	}
}

func (h *Hub) addClient(client *Client) {
	if _, ok := h.clients[client.UUID]; !ok {
		h.clients[client.UUID] = make(map[*Client]bool)
	}
	h.clients[client.UUID][client] = true
}
