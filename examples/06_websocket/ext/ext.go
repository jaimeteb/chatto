package main

import (
	"github.com/jaimeteb/chatto/query"

	"github.com/jaimeteb/chatto/internal/channels/message"

	"github.com/jaimeteb/chatto/extensions"
)

func helloUniverseFunc(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
	return &extensions.ExecuteExtensionResponse{
		FSM: req.FSM,
		Response: message.Response{
			Answers: []query.Answer{{
				Text:  "Hello Universe",
				Image: "https://i.imgur.com/pPdjh6x.jpg",
			}},
			ReplyOpts: req.Request.ReplyOpts,
			Channel:   req.Request.Channel,
		},
	}
}

var registeredExtensions = extensions.Registered{
	"any": helloUniverseFunc,
}

func main() {
	websocketClient := extensions.NewWebSocketClient("127.0.0.1:4770", registeredExtensions)

	websocketClient.StartWebsocketClient()
}
