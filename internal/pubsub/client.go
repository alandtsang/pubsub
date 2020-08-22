// client package saves the websocket connection and
// is responsible for read and write messages.
package pubsub

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client contains the websocket connection and a randomly generated id.
type Client struct {
	// The uuid of client
	id string

	// IP from websocket connection
	ip string

	// Websocket connection
	conn *websocket.Conn

	// Send buffer channel
	send chan []byte

	// Pubsub pointer
	ps *Pubsub

	// Whether the websocket connection is closed
	closed bool
}

// New returns a new websocket client.
func NewClient(conn *websocket.Conn, ip string, ps *Pubsub) *Client {
	return &Client{
		id:   uuid.Must(uuid.NewV4()).String(),
		ip:   ip,
		conn: conn,
		send: make(chan []byte, 100),
		ps:   ps,
	}
}

// GetID returns the id of the Client.
func (c *Client) GetID() string {
	return c.id
}

// GetIP returns the ip of websocket connection.
func (c *Client) GetIP() string {
	return c.ip
}

// Read receives message from websocket connection, and process the message.
func (c *Client) Read() {
	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, p, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("client %s read message failed, %v", c.id, err)
			return
		}

		log.Printf("client %s recv: %s", c.id, string(p))
		p = bytes.TrimSpace(bytes.Replace(p, newline, space, -1))
		process(c, p)
	}
}

// Write writes sends ping message and received message to websocket connection.
func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Printf("client %s is disconnected", c.id)
				return
			}
			_ = c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ping client %s failed\n", c.id)
				c.ps.RemoveClient(c)
				return
			}
		}
	}
}

// Close closes websocket connection.
func (c *Client) Close() {
	if !c.closed {
		_ = c.conn.Close()
		c.closed = true
	}
}

func process(cli *Client, p []byte) {
	var msg Message
	if err := json.Unmarshal(p, &msg); err != nil {
		cli.send <- []byte("invalid message")
		return
	}

	switch msg.Action {
	case actionTypeSubscribe:
		cli.ps.Subscribe(msg.Topic, cli)
	case actionTypePublish:
		if len(msg.Msg) == 0 {
			log.Println("invalid message content")
			return
		}

		cli.ps.Publish(msg.Topic, msg.Msg)
	default:
		cli.send <- []byte("invalid action")
	}
}
