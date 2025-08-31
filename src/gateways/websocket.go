package gateways

import (
	"fmt"

	"github.com/gofiber/contrib/websocket"
)

func (h *HTTPGateway) AskChatWs(wsconn *websocket.Conn) {

	fmt.Println("WebSocket connected. Waiting for messages...")

	defer func() {
		fmt.Println("WebSocket connection closed.")
		wsconn.Close()
	}()

	for {
		// perform read question and answer that question
	}

}
