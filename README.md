[![Build Status](https://travis-ci.com/jaimeteb/chatto.svg?branch=master)](https://travis-ci.com/jaimeteb/chatto)
[![codecov](https://codecov.io/gh/jaimeteb/chatto/branch/master/graph/badge.svg)](https://codecov.io/gh/jaimeteb/chatto)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaimeteb/chatto)](https://goreportcard.com/report/github.com/jaimeteb/chatto)
[![GoDoc](https://godoc.org/github.com/jaimeteb/chatto?status.svg)](https://godoc.org/github.com/jaimeteb/chatto)
[![Docker Image Version (latest by date)](https://img.shields.io/docker/v/jaimeteb/chatto?color=teal&sort=date)](https://hub.docker.com/repository/docker/jaimeteb/chatto)

---
# chatto

<p align="center">
<img src="https://user-images.githubusercontent.com/17936011/89082867-e3c0d300-d354-11ea-9def-008c403a4497.png" alt="botto" width="150"/>
</p>

Simple chatbot framework written in Go, with configurations in YAML. The aim of this project is to create very simple text-based chatbots using a few configuration files. 

The inspiration for this project originally came from [Flottbot](https://github.com/target/flottbot) and my experience using [Rasa](https://github.com/RasaHQ/rasa).

## Contents

* [Installation](#install)
* [Usage](#usage)
  * [How does it work?](#how)
* [Classifier](#clf)
  * [Pipeline](#pipeline)
* [Finite State Machine](#fsm)
  * [*Any*](#any)
  * [Default Messages](#defaults)
  * [Slots](#slots)
* [Extensions](#ext)
* [Bot Configuration](#botconfig)
  * [Extensions](#botconfigext)
  * [Store](#store)
  * [Channels](#channels)
    * [Telegram](#telegram)  
    * [Twilio](#twilio)
* [Endpoints](#endpoints)
  * [HTTP Endpoint](#http)  
  * [CLI](#cli)
  * [Prediction](#predict)
* [Examples](#examples)  

<a name="install"></a>
## Installation

Run `go get -u github.com/jaimeteb/chatto`.

<a name="usage"></a>
## Usage

Run `chatto` in the directory where your YAML files are located, or specify a path to them with the `-path` flag:

```
chatto -path ./your/data
```

To run on Docker, use:
```bash
docker run \
  -p 4770:4770 \
  -e CHATTO_DATA=./your/data \
  -v $PWD/your/data:/chatto/data \
  jaimeteb/chatto
```

<a name="how"></a>
### How does it work?

Chatto combines the consistency of a finite-state-machine with the flexibility of machine learning. It has three main components: the classifier, the finite-state-machine and the extensions.

A very basic directory structure for chatto would be the following:

```
data
├── clf.yml
└── fsm.yml
```

<a name="clf"></a>
## Classifier

Currently, chatto uses a [Naïve-Bayes classifier](github.com/navossoc/bayesian) to take the user input and decide a command to execute on the finite-state-machine. The training text for the classifier is provided in the **clf.yml** file:

```yaml
classification:
  - command: "turn_on"
    texts:
      - "turn on"
      - "on"

  - command: "turn_off"
    texts:
      - "turn off"
      - "off"
```

Under **classification** you can list the commands and their respective training data under **texts**.

<a name="pipeline"></a>
### Pipeline

You can optionally configure the pipeline steps (currently: removal of symbols, conversion into lowercase and classification threshold) adding the *pipeline* object to the **clf.yml** file:
```yaml
pipeline:
  remove_symbols: true
  lower: true
  threshold: 0.3
```

<a name="fsm"></a>
## Finite State Machine

The FSM (finite-state-machine) is based on the one shown in [this article](https://levelup.gitconnected.com/implement-a-finite-state-machine-in-golang-f0438b6bc0a8). The states, commands, default messages and transitions are described in the **fsm.yml** file:

```yaml
states:
  - "off"
  - "on"
  
commands:
  - "turn_on"
  - "turn_off"

functions:
 - transition:
      from: "off"
      into: "on"
    command: "turn_on"
    message: "Turning on."

  - transition:
      from: "on"
      into: "off"
    command: "turn_off"
    message:
      - "Turning off."
      - "❌"

defaults:
  unknown: "Can't do that."
```

Under **functions** you can list the transitions available for the FSM. The object **transition** describes the states of the transition (**from** one state **into** another) if **command** is executed; **message** is the message (or messages) to send to the user.

<a name="any"></a>
### *Any*

The special state **any** can help you go from any state into another, if the command is executed. You don't have to declare the **any** state in the states list.

Also, the special command **any** is used to transition between two states, regardless of the command predicted. This command doesn't have to be declared.

<a name="defaults"></a>
### Default Messages

In the **fsm.yml** file, the *defaults* section is used to set the messages that will be returned when the following events happen:

- **unknown**: The current state does not transition into another one with the predicted command.
- **unsure**: The command prediction confidence was below the threshold.
- **error**: An error ocurred during the execution of an extension.

Here's an example of default messages:

```yaml
defaults:
  unknown: "Can't do that transition."
  unsure: "Sorry, I didn't understand that."
  error: "An error ocurred."
```

<a name="slots"></a>
### Slots

You can save information from the user's input by using slots. In the **fsm.yml**, slots are declared as such:

```yaml
  - transition:
      from: ask_name
      into: ask_age
    command: say_name
    slot:
      name: name
      mode: whole_text
    message: "How old are you?"
```

In this example, in the transition from **ask_name** to **ask_age**, when **say_name** is executed, a slot called **name** will be saved, in other words, the user's message is stored in memory.

At the time, only **whole_text** mode is supported, which saves the entire input in the slot.

<a name="ext"></a>
## Extensions

The extensions in chatto are pieces of Go code that can be executed instead of messages, and can also alter the state of the conversation. In the **fsm.yml** file, the extensions' names must begin by **"ext_"**.
Extensions are executed as services in a separate Go file, for example **ext.go**.

The format for a chatto extension is as follows:

```go
package main

import (
	"log"

	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(req *fsm.Request) (res *fsm.Response) {
	return &fsm.Response{
		FSM: req.FSM,
		Res: "Hello Universe",
	}
}

var myExtMap = fsm.ExtensionMap{
	"ext_any": greetFunc,
}

func main() {
	log.Fatal(fsm.ServeExtensionREST(myExtMap))
}
```

You must use either `ServeExtensionRPC` or `ServeExtensionREST` in the main function in order to run the extension server and pass your own **fsm.ExtensionMap**, which maps the extension names to their respective functions.

There are currently two ways to serve the extensions:

- **RPC**: By using `fsm.ServeExtensionRPC(fsm.ExtensionMap)`
- **REST**: By using `fsm.ServeExtensionREST(fsm.ExtensionMap)`

When running the extensions, use the flag ```-port``` to specify a service port (extensions will use port 8770 by default).

The extension functions must have the ```func(*fsm.Request) *fsm.Response``` signature, where:
* Request contains:
  * The current FSM
  * The requested extension
  * The input text from the user
  * The Domain (*fsm.yml* data)
* Response must contain:
  * The resulting FSM
  * The message to be sent to the user

In this example, **ext_any** simply returns "Hello Universe" and does not modify the current FSM.

<a name="botconfig"></a>
## Bot Configuration

The **bot.yml** file is used to configure the name of the bot, how and where the extensions will be consumed, and how will the FSMs will be stored.

```yaml
bot_name: "test_bot"
extensions:
  type: REST
  url: http://localhost:8770
store:
  type: REDIS
  host: localhost
  password: pass
```

<a name="botconfigext"></a>
### Extensions

To configure the extensions, the following parameters are required for RPC and REST types respectively:

- For type **RPC**:
  - Host
  - Port
- For type **REST**:
  - URL

<a name="store"></a>
### Store

The FSMs for the bot can be stored locally (default type **CACHE**) or in Redis. In order to use Redis, provide the following values, as shown in the example above:

- For type **REDIS**:
  - Host
  - Password

You can leave the values empty and set them with environment variables, for example:

```yaml
extensions:
  type: RPC
  host: 
  port: 
store:
  type: REDIS
  host: 
  password: 
```

And set the environment variables:

```
EXTENSIONS_HOST=localhost
EXTENSIONS_PORT=8770
STORE_HOST=localhost
STORE_PASSWORD=pass
```

<a name="channels"></a>
### Channels

In the **chn.yml** you can insert the credentials for a Telegram Bot and/or a Twilio phone number.

```yaml
telegram:
  bot_key: MY_BOT_KEY
twilio:
  account_sid: MY_ACCOUNT_SID
  auth_token: MY_AUTH_TOKEN
  number: MY_NUMBER
```

<a name="telegram"></a>
#### Telegram

You can connect your chatto bot to [Telegram](https://core.telegram.org/bots) by providing yout Telegram Bot Key, either directly in the **chn.yml** file or by setting the **TELEGRAM_BOT_KEY** environment variable.

You must set the bot's webhook to the ***/endpoints/telegram*** endpoint in order to receive messages.

<a name="twilio"></a>
#### Twilio

Similarly, connect your bot to [Twilio](https://www.twilio.com/messaging-api) by adding your credentials to the file or by setting the corresponding environment variables (**TWILIO_ACCOUNT_SID**, etc.)

You must set the webhooks to the ***/endpoints/twilio*** endpoint in order to receive messages.

<a name="endpoints"></a>
## Endpoints

Apart from the [channels](#channels), once your chatto bot is running, you can interact with it via the HTTP endpoints.

<a name="http"></a>
### HTTP Endpoint

Chatto will run on port 4770 bu default. You can specify a different one with the `-port` flag.

Send a *POST* request to */endpoints/rest* with the following body structure:

```json
{
    "sender": "foo",
    "text": "bar"
}
```

Example with cURL:

```bash
curl --request POST 'http://localhost:4770/endpoints/rest' \
--header 'Content-Type: application/json' \
--data-raw '{
    "sender": "foo",
    "text": "bar"
}'
```

The bot will respond as such:

```json
{
    "sender": "botto",
    "text": "some answer"
}
```

<a name="cli"></a>
### CLI

To enable the CLI mode, run chatto using the `-cli` flag.

This will launch a command line interface where you can send and receive messages from your bot. This is a useful mode when debugging.

<a name="predict"></a>
### Prediction

You can test the command predictions with a *POST* request to the */predict* endpoint with the following body structure:

```json
{
    "text": "good"
}
```

Example with cURL:

```bash
curl --request POST 'http://localhost:4770/predict' \
--header 'Content-Type: application/json' \
--data-raw '{
    "text": "foo"
}'
```

The resulting prediction will look like this:

```json
{
    "original": "foo",
    "predicted": "good",
    "probability": 0.3274336283185841
}
```

<a name="examples"></a>
## Examples

I have provided some config files under *examples*. Clone the repository and run `chatto` with the `-path` of your desired example to test them out (for the ones that use extensions, run their respective extensions first).

1. [**Mood Bot**](/examples/01_moodbot) - A chatto version of [Rasa's Mood Bot](https://github.com/RasaHQ/rasa/tree/master/examples/moodbot) Greet the bot to start the conversation.
2. [**Engineering Flowchart**](/examples/02_repair) - Tell the bot you want to repair something.
3. [**Pokemon**](/examples/03_pokemon) - Search for Pokémon by name or number.
4. [**Trivia Quiz**](/examples/04_trivia) - Type *start* to take a quick trivia quiz.
