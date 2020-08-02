package ws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/alandtsang/pubsub/internal/client"
	"github.com/alandtsang/pubsub/internal/pubsub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ps *pubsub.Pubsub

func init() {
	ps = pubsub.New()
}

// WebsocketHandler is responsible for handling the websocket connection.
func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cli := client.New(conn)
	fmt.Printf("New client %s connected\n", cli.GetID())
	_ = cli.WriteTextMessage(fmt.Sprintf("Hi client %s", cli.GetID()))

	for {
		ps.HandleMessage(cli)
	}
}
