package pubsub

import (
	"encoding/json"
	"fmt"

	"github.com/alandtsang/pubsub/internal/client"
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

// Pubsub is responsible for handling publish and subscribe.
type Pubsub struct {
	// saving clients that subscribe the topic.
	subscriptions map[string][]*client.Client
}

// New return a new Pubsub structure pointer.
func New() *Pubsub {
	return &Pubsub{
		subscriptions: make(map[string][]*client.Client),
	}
}

// HandleMessage handles pub and sub messages.
func (ps *Pubsub) HandleMessage(cli *client.Client) {
	p, err := cli.ReadMessage()
	if err != nil {
		return
	}

	var msg Message
	if err = json.Unmarshal(p, &msg); err != nil {
		fmt.Println("ReadMessage unmarshal message failed,", err)
		_ = cli.WriteTextMessage("invalid message")
		return
	}

	fmt.Printf("action: %v, topic: %s, msg: %s\n", msg.Action, msg.Topic, string(msg.Msg))

	switch msg.Action {
	case actionTypeSubscribe:
		ps.subscribe(msg.Topic, cli)
	case actionTypePublish:
		if len(msg.Msg) == 0 {
			_ = cli.WriteTextMessage("invalid message")
			return
		}
		ps.publish(msg.Topic, msg.Msg)
	default:
		_ = cli.WriteTextMessage("invalid action")
	}
}

func (ps *Pubsub) publish(topic string, message interface{}) {
	clients, ok := ps.subscriptions[topic]
	if !ok {
		fmt.Printf("invalid topic %s\n", topic)
		return
	}

	msg, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("publish json marshal message failed, %v", err)
		return
	}

	for _, cli := range clients {
		_ = cli.WriteTextMessage(string(msg))
	}
}

func (ps *Pubsub) subscribe(topic string, cli *client.Client) {
	ps.subscriptions[topic] = append(ps.subscriptions[topic], cli)
}
