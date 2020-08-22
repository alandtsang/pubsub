package main

import (
	"encoding/json"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

// Message is the structure of subscribing and publishing messages.
type Message struct {
	// Action type, publish or subscribe
	Action string `json:"action"`

	// The topic used to subscribe or publish messages
	Topic string `json:"topic"`

	// Message to send
	Msg string `json:"msg"`
}

type client struct {
	addr  string
	path  string
	query string
	conn  *websocket.Conn
	buf   chan []byte
}

func newClient(addr, path, query string) *client {
	return &client{
		addr:  addr,
		path:  path,
		query: query,
		buf:   make(chan []byte),
	}
}

func (c *client) dial() {
	u := url.URL{Scheme: "ws", Host: c.addr, Path: c.path, RawQuery: c.query}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c.conn = conn
}

func (c *client) receive() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}

func (c *client) subscribe(topic string) {
	message := &Message{
		Action: "subscribe",
		Topic:  topic,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("subscribe json marshal failed, %v", err)
		return
	}
	c.buf <- data
}

func (c *client) publish(topic, msg string) {
	message := &Message{
		Action: "publish",
		Topic:  topic,
		Msg:    msg,
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("publish json marshal failed, %v", err)
		return
	}
	c.buf <- data
}

func (c *client) write() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.buf:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func (c *client) close() {
	_ = c.conn.Close()
}

func main() {
	addr := "localhost:9999"
	path := "/ws"
	query := ""

	client := newClient(addr, path, query)
	client.dial()
	defer client.close()

	go client.receive()
	go client.write()

	topic := "chat"
	client.subscribe(topic)
	client.publish(topic, "This is a test message")

	time.Sleep(10 * time.Second)
}
