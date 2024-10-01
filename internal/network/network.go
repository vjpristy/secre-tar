package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Connection struct {
	*websocket.Conn
}

func NewServer(addr string) *http.Server {
	return &http.Server{
		Addr: addr,
	}
}

func HandleConnections(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &Connection{conn}, nil
}

func (c *Connection) ReadMessage() ([]byte, error) {
	_, message, err := c.Conn.ReadMessage()
	return message, err
}

func (c *Connection) WriteMessage(message []byte) error {
	return c.Conn.WriteMessage(websocket.TextMessage, message)
}

func Dial(urlStr string) (*Connection, *http.Response, error) {
	c, resp, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		return nil, resp, err
	}
	return &Connection{c}, resp, nil
}
