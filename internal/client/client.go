package client

import (
	"fmt"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

// Client contains the websocket connection and a randomly generated id.
type Client struct {
	id   string
	conn *websocket.Conn
}

// New returns a new websocket client.
func New(conn *websocket.Conn) *Client {
	return &Client{
		id:   uuid.Must(uuid.NewV4()).String(),
		conn: conn,
	}
}

// GetID returns the id of the Client.
func (c *Client) GetID() string {
	return c.id
}

// ReadMessage receives message from websocket connection,
// and return the read data.
func (c *Client) ReadMessage() ([]byte, error) {
	messageType, p, err := c.conn.ReadMessage()
	if err != nil {
		fmt.Printf("client %s read message failed, %v\n", c.id, err)
		return nil, err
	}
	fmt.Println("message type", messageType)

	return p, nil
}

// WriteTextMessage sends message to the websocket connection.
func (c *Client) WriteTextMessage(p string) error {
	return c.conn.WriteMessage(websocket.TextMessage, []byte(p))
}
