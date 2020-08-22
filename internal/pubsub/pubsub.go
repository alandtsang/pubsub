// Package pubsub stores clients connection, and responsible
// for subscribing and publishing messages
package pubsub

import (
	"encoding/json"
	"log"
	"sync"
)

type actionType string

const (
	actionTypeSubscribe actionType = "subscribe"
	actionTypePublish   actionType = "publish"
)

// Message is the structure of subscribing and publishing messages.
type Message struct {
	// Action type, publish or subscribe
	Action actionType `json:"action"`

	// The topic used to subscribe or publish messages
	Topic string `json:"topic"`

	// Message to send
	Msg json.RawMessage `json:"msg"`
}

// Pubsub saves the connected clients, and responsible for handling publish and subscribe.
type Pubsub struct {
	// The mutex to protect connections
	mutex sync.RWMutex

	// Registered clients.
	clients map[*Client]map[string]struct{}
}

// New return a new Pubsub structure pointer.
func NewPubSub() *Pubsub {
	return &Pubsub{
		clients: make(map[*Client]map[string]struct{}),
	}
}

// AddClient adds new client to map.
func (ps *Pubsub) AddClient(cli *Client) {
	ps.mutex.Lock()
	ps.clients[cli] = make(map[string]struct{})
	ps.mutex.Unlock()
	log.Printf("New client %s from %s connected", cli.GetID(), cli.GetIP())
}

// RemoveClient removes client from the map.
func (ps *Pubsub) RemoveClient(cli *Client) {
	ps.mutex.Lock()
	if _, ok := ps.clients[cli]; ok {
		delete(ps.clients, cli)
		cli.Close()
		log.Printf("Remove client %s", cli.GetID())
	}
	ps.mutex.Unlock()
}

// Publish publishes msg to all clients subscribed to topic.
func (ps *Pubsub) Publish(topic string, msg []byte) {
	for cli, topics := range ps.clients {
		if _, ok := topics[topic]; ok {
			cli.send <- msg
		}
	}
}

// Subscribe records the topic subscribed by the client.
func (ps *Pubsub) Subscribe(topic string, cli *Client) {
	ps.clients[cli][topic] = struct{}{}
}
