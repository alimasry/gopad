package ws_test

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/alimasry/gopad/internal/pkg/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func TestWebSocketHandler(t *testing.T) {
	// Set up Gin router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ws", ws.ServeWs)

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Create WebSocket connection to test server
	url := "ws" + server.URL[len("http"):len(server.URL)] + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("Dial failed: %v", err)
	}
	defer ws.Close()
	fmt.Println(url)
}
