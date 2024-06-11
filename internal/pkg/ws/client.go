package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/alimasry/gopad/internal/services/editor"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	syncPeriod = 500 * time.Millisecond
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientList map[string]map[*Client]bool

type Client struct {
	ReplicaId         string
	conn              *websocket.Conn
	hub               *Hub
	send              chan *Event
	documentUUID      string
	lastSyncedVersion int
}

// creates a new websocket client
func NewClient(documentUUID string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ReplicaId:    uuid.NewString(),
		conn:         conn,
		hub:          hub,
		documentUUID: documentUUID,
		send:         make(chan *Event, 256),
	}
}

// receives messages from the client
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		c.routeMessage(message)
	}
}

func (c *Client) routeMessage(message []byte) {
	var event Event

	if err := json.Unmarshal(message, &event); err != nil {
		log.Printf("error: %v", err)
	}

	event.client = c

	routeEvent(event)
}

// sends messages to the client
func (c *Client) writePump() {
	pingTicker := time.NewTicker(pingPeriod)
	syncTicker := time.NewTicker(syncPeriod)

	defer func() {
		pingTicker.Stop()
		syncTicker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			message, err := json.Marshal(event)
			if err != nil {
				log.Printf("error : %v", err)
			}

			err = c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("error : %v", err)
			}
		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case <-syncTicker.C:
			c.syncIfChanged()
		}

	}
}

func (c *Client) syncIfChanged() {
	document := editor.GetDocumentFromCache(c.documentUUID)

	if c.lastSyncedVersion == document.Version {
		return
	}

	syncData, err := json.Marshal(document)
	if err != nil {
		log.Printf("error: %v", err)
	}

	c.lastSyncedVersion = document.Version

	c.send <- &Event{
		Command: SyncEvent,
		Data:    syncData,
		client:  c,
	}
}

// handles the ws endpoint
func ServeWs(c *gin.Context) {
	hub := GetHubInstance()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	documentUUID := c.Param("document_uuid")

	client := NewClient(documentUUID, conn, hub)
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
