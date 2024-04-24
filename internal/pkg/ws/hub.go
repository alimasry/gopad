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

// initialize the hub instance
func initializeHub() {
	once.Do(func() {
		hubInstance = &Hub{
			clients:    make(ClientList),
			broadcast:  make(chan *Event),
			register:   make(chan *Client),
			unregister: make(chan *Client),
		}
	})
}

// get the hub instance and initialize it if it's not
func GetHubInstance() *Hub {
	initializeHub()
	return hubInstance
}

// runs a loop that registers / unregisters clients as well as handle broadcast events
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.addClient(client)
		case client := <-h.unregister:
			h.removeClient(client)
		case event := <-h.broadcast:
			for client := range h.clients[event.client.documentUUID] {
				select {
				case client.send <- event:
				default:
					h.removeClient(client)
				}
			}
		}
	}
}

// remove client from the hub
func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client.documentUUID][client]; ok {
		delete(h.clients[client.documentUUID], client)
		close(client.send)

		if len(h.clients) == 0 {
			delete(h.clients, client.documentUUID)
		}
	}
}

// add client to the hub
func (h *Hub) addClient(client *Client) {
	if _, ok := h.clients[client.documentUUID]; !ok {
		h.clients[client.documentUUID] = make(map[*Client]bool)
	}
	h.clients[client.documentUUID][client] = true
}
