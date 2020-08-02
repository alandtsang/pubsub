# pubsub

Publish and subscribe demo with websocket in Go

## Demo

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

`Message` is the structure of subscribing and publishing messages.

```go
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
```

`HandleMessage` handles pub and sub messages.

```go
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
```

## Usage

Execute the following in the terminal：

```
go run cmd/main.go
```

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

## License

Please refer to [LICENSE](https://github.com/alandtsang/pubsub/blob/master/LICENSE) file.