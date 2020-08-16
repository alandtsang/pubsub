// Package route contains homepage and websocket routing.
package route

import (
	"fmt"
	"log"
	"net/http"

	ws "github.com/alandtsang/pubsub/pkg/websocket"
)

func SetRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", ws.WebsocketHandler)

	fmt.Println("Listen on localhost:9999")
	log.Println(http.ListenAndServe(":9999", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web")
}
