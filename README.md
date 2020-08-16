# pubsub

Publish and subscribe demo with websocket in Go

## Demo

### Html

Initialize websocket in static `index.html`:

```javascript
<script type="text/javascript">
    var ws = new WebSocket("ws://localhost:9999/ws");
    ws.onopen = () => {
        console.log("Connected");
    }

    ws.onclose = () => {
        console.log("Disconnected");
    }

    ws.onmessage = (msg) => {
        console.log("Server message:", msg.data)
    }

    window.ws = ws
</script>
```

### Message

`Message` is the structure of subscribing and publishing messages.

```go
// Message is the structure of subscribing and publishing messages.
type Message struct {
    // Action type, publish or subscribe
    Action actionType `json:"action"`

    // The topic used to subscribe or publish messages
    Topic string `json:"topic"`

    // Message to send
    Msg json.RawMessage `json:"msg"`
}
```

### Client

`Client` contains the websocket connection and a randomly generated id.

```go
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
```

### Pubsub

`Pubsub` saves the connected clients, and responsible for handling publish and subscribe.

```go
type Pubsub struct {
    // The mutex to protect connections
    mutex sync.RWMutex

    // Registered clients.
    clients map[*Client]map[string]struct{}
}
```

## Usage

Execute the following in the terminal：

```
go run cmd/main.go
```

### browser

You can enter in the browser's Console:

```javascript
ws.send('{"action":"subscribe","topic":"abc"}')
```

output:

```
Server message: "This is a broadcast"
```

```javascript
ws.send('{"action":"publish","topic":"abc","msg":"This is a broadcast"}')
```

output:

```
Server message: "This is a broadcast"
```

### wscat

```
wscat -c "ws://localhost:9999/ws"
Connected (press CTRL+C to quit)
> {"action":"subscribe","topic":"abc"}
> {"action":"publish","topic":"abc","msg":"This is a broadcast"}
< "This is a broadcast"
```

## License

Please refer to [LICENSE](https://github.com/alandtsang/pubsub/blob/master/LICENSE) file.