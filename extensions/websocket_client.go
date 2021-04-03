package extensions

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketClient is an extension client that receives execute request events
type WebSocketClient struct {
	Address              string
	RegisteredExtensions Registered
}

func NewWebSocketClient(address string, registeredExtensions Registered) *WebSocketClient {
	return &WebSocketClient{
		Address:              address,
		RegisteredExtensions: registeredExtensions,
	}
}

func (s *WebSocketClient) StartWebsocketClient() {
	u := url.URL{Scheme: "ws", Host: s.Address, Path: "/bot/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	defer func() {
		err = c.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}

		log.Printf("recv: %s", message)

		time.Sleep(5 * time.Second)
	}
}
