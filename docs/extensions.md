# Extensions

The extensions in Chatto are pieces of code that can be executed instead of messages, and can also alter the state of the conversation. In the **fsm.yml** file, the extensions are contained in the `extension` field, under `transitions`.

Extensions are executed as services, and can be written in Go, using the [`chatto/extensions`](https://godoc.org/github.com/jaimeteb/chatto/extensions) and [`chatto/query`](https://godoc.org/github.com/jaimeteb/chatto/query) packages, or they can be written in any language, as long as the services are compatible.

## Go

In Golang, the format for a Chatto extension server is as follows:

```go
package main

import (
	"log"

	"github.com/jaimeteb/chatto/extensions"
	"github.com/jaimeteb/chatto/query"
)

// GreetFunc returns the message "Hello Universe" and an image
// and does not modify the Finite State Machine (FSM)
func GreetFunc(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
	return &extensions.ExecuteExtensionResponse{
		FSM: req.FSM,
		Answers: []query.Answer{{
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		}},
	}
}

// registeredExtensions maps the name "any" to the GreetFunc extension
var registeredExtensions = extensions.RegisteredExtensions{
	"any": GreetFunc,
}

func main() {
	// Run the extensions via REST
	log.Fatal(extensions.ServeREST(registeredExtensions))
}
```

You must use either `ServeRPC` or `ServeREST` in the main function in order to run the extension server and pass your own [`extensions.RegisteredExtensions`](https://godoc.org/github.com/jaimeteb/chatto/extensions#RegisteredExtensions), which maps the extension names to their respective functions.

There are currently two ways to serve the extensions:

- **RPC**: By using [`extensions.ServeRPC`](https://godoc.org/github.com/jaimeteb/chatto/extensions#ServeRPC)
- **REST**: By using [`extensions.ServeREST`](https://godoc.org/github.com/jaimeteb/chatto/extensions#ServeREST)

When running the extensions, use the flag `-port` to specify a service port (extensions will use port `8770` by default).

The extension functions must have the signature:

```go
func(*extensions.ExecuteExtensionRequest) *extensions.ExecuteExtensionResponse
```

Where:

* [`ExecuteExtensionRequest`](https://godoc.org/github.com/jaimeteb/chatto/extensions#ExecuteExtensionRequest) contains:
	* The current FSM
	* The channel that received the request
	* The requested extension
	* The input question (the sender and the text)
	* The Domain (*fsm.yml* data)
* [`ExecuteExtensionResponse`](https://godoc.org/github.com/jaimeteb/chatto/extensions#ExecuteExtensionResponse) must contain:
	* The resulting FSM
	* The answers (messages) to be sent to the user

In this example, the extension **any** simply returns "Hello Universe" and an image, and does not modify the current FSM.

## Other languages

Since extensions are services, they can be written in any language. Here is an example in Python, that is equivalent to the one shown above in Go.

### Python example

```python
from flask import Flask, Response, request, jsonify

app = Flask(__name__)

def greet_func(data: dict) -> dict:
	return jsonify({
		"fsm": data.get("fsm"),
		"answers": [
			{
				"text": "Hello Universe",
				"image": "https://i.imgur.com/pPdjh6x.jpg",
			}
		]
	})

registered_extensions = {
    "<any": greet_func,
}

@app.route("/extensions", methods=["GET"])
def get_all_funcs():
    return jsonify(list(registered_extensions.keys()))

@app.route("/extension", methods=["POST"])
def get_func():
    data = request.get_json()
    req = data.get("extension")
    f = registered_extensions.get(req)
    if not f:
        return Response(status=400)
    return f(data)

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8770)
```

In this case, the Flask app emulates the function of [`ServeREST`](https://godoc.org/github.com/jaimeteb/chatto/extensions#ServeREST), while `greet_func` and `registered_extensions` correspond to `GreetFunc` and `registeredExtensions` respectively, from the Go example.

### Extension REST

An extension REST service must implement these routes:

- `GET /extensions`

	This route should return an array with the names of the registered extensions.

	Example response:
	```json
	[
		"val_ans_1",
		"val_ans_2",
		"score"
	]
	```

- `POST /extension`

	This route should return the resulting Finite State Machine object after the extension's execution, along with the answers.

	Example request body:

	```json
	{
		"fsm": {
			"state": 2,
			"slots": {
				"answer_1": "3"
			}
		},
		"extension": "val_ans_1",
		"channel": "rest",
		"question": {
			"sender": "cli",
			"text": "2"
		},
		"domain": {
			"state_table": {
				"any": -1,
				"initial": 0,
				"question_1": 1,
				"question_2": 2,
				"question_3": 3
			},
			"command_list": [
				"start",
				"end"
			],
			"default_messages": {
				"unknown": "Not sure I understood, try again please.",
				"unsure": "Not sure I understood, try again please.",
				"error": "I'm sorry, there was an error."
			}
		}
	}
	```

	Example response:

	```json
	{
		"fsm": {
			"state": 2,
			"slots": {
				"answer_1": "3"
			}
		},
		"answers": [
			{
				"text": "Select one of the options."
			}
		]
	}
	```

## Answers

The answers returned from the extensions follow the same rules as [the **fsm.yml** messages](/finitestatemachine/#messages). In Go, you can use the helper [`query.Answers`](https://godoc.org/github.com/jaimeteb/chatto/query#Answers) function to create answers from [`query.Answer`](https://godoc.org/github.com/jaimeteb/chatto/query#Answer), strings or maps.

```go
func GreetFunc(req *extensions.ExecuteExtensionRequest) (res *extensions.ExecuteExtensionResponse) {
	return &extensions.ExecuteExtensionResponse{
		FSM: req.FSM,
		// Answers: query.Answers("Hello Universe"),		// This is a simple string answer
		// Answers: query.Answers("Hello", "Universe"),		// This is a slice of answers
		Answers: query.Answers(query.Answer{				// This is a text/image answer
			Text:  "Hello Universe",
			Image: "https://i.imgur.com/pPdjh6x.jpg",
		})
	}
}
```

In REST Extensions in other languages, answers must meet the corresponfing JSON notation. The following JSON are valid `answers` messages.

```json
{
	"answers": [
		{
			"text": "A simple answer message"
		}
	]
}

{
	"answers": [
		{
			"text": "A list"
		},
		{
			"text": "of answers"
		}
	]
}

{
	"answers": [
		{
			"text": "An image will be attached"
		},
		{
			"text": "to one of these answers",
			"image": "https://i.imgur.com/8MU0IUT.jpeg"
		}
	]
}
```
