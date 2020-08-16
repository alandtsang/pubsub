// Package ws is the handler of websocket, create a new client.
package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/alandtsang/pubsub/internal/pubsub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ps *pubsub.Pubsub

func init() {
	ps = pubsub.NewPubSub()
}

// WebsocketHandler is responsible for handling the websocket connection.
func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cli := pubsub.NewClient(conn, r.RemoteAddr, ps)
	ps.AddClient(cli)
	go cli.Read()
	go cli.Write()
}
