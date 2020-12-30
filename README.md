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

[1. Installation](#install)  
[2. How does it work?](#how)  
[3. Classifier](#clf)  
[4. Finite State Machine](#fsm)  
[5. Extensions](#ext)  
[6. Slots](#slots)  
[7. Redis](#redis)  
[8. Pipeline](#pipeline)  
[9. HTTP Endpoint](#endpoint)  
[10. CLI](#cli)  
[11. Telegram](#telegram)  
[12. Examples](#examples)  

<a name="install"></a>
## 1. Installation

Run ```go get -u github.com/jaimeteb/chatto```.

<a name="how"></a>
## 2. How does it work?

Chatto combines the consistency of a finite-state-machine with the flexibility of machine learning. It has three main components: the classifier, the finite-stete-machine and the extensions.

<a name="clf"></a>
## 3. Classifier

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

<a name="fsm"></a>
## 4. Finite State Machine

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
  unsure: "???"
```

Under **functions** you can list the transitions available for the FSM. The object **transition** describes the states of the transition (**from** one state **into** another) if **command** is executed; **message** is the message (or messages) to send to the user.

The special state **any** can help you go from any state into another, if the command is executed. You don't have to declare the **any** state in the states list.

<a name="ext"></a>
## 5. Extensions

The extensions in chatto are pieces of code that can be executed instead of messages. The extensions names must begin by **"ext_"** and they must be placed in the **ext/ext.go** file. The format for a chatto extension is as follows:

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
	if err := fsm.ServeExtension(myExtMap); err != nil {
		log.Fatalln(err)
	}
}
```

You must use the ```fsm.ServeExtension(fsm.ExtensionMap)``` in the main function in order to run the extension server and pass your own **fsm.ExtensionMap**, which maps the extension names to their respective functions.

The extension server runs on port ```42586``` by default but you can specify it with the **EXTENSION_PORT** environment variable. Furthermore, you can run the extension server elsewhere, in which case you have to ser the **EXTENSION_HOST** environment variable.

The extension functions must have the ```func(*fsm.Request) *fsm.Response interface{}``` signature, where:
* Request contains:
  * The current FSM
  * The requested extension
  * The input text from the user
  * The Domain
* Response must contain:
  * The resulting FSM
  * The message to be sent to the user

In this example, **ext_any** simply returns "Hello Universe" and does not modify the current FSM.

<a name="slots"></a>
## 6. Slots

You can save information from the user's input by using slots:

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

<a name="redis"></a>
## 7. Redis

You can store the FSMs in memory or in Redis. In order to use the Redis Store, set the **REDIS_HOST** and **REDIS_PASS** environment variables.

<a name="pipeline"></a>
## 8. Pieline

You can optionally configure the pipeline steps (removal of symbols, conversion into lowercase and classification threshold) using the **pl.yml** file:
```yaml
remove_symbols: true
lower: true
threshold: 0.3
```

<a name="endpoint"></a>
## 9. HTTP Endpoint

To enable the HTTP endpoint, simply run ```chatto``` on the same directory as your **clf.yml** and **fsm.yml** files, or specify a path to them with the ```--path``` flag. A service will run on port 4770 of your localhost.

Send a *POST* request to */endpoint* with the following body structure:

```json
{
    "sender": "foo",
    "text": "bar"
}
```

The bot will respond as such:

```json
{
    "sender": "botto",
    "text": "some answer"
}
```

<a name="cli"></a>
## 10. CLI

Alternatively, run chatto on a command line interface using the ```--cli``` flag.

<a name="telegram"></a>
## 11. Telegram

You can connect your chatto bot to [Telegram](https://core.telegram.org/bots) by setting the **TELEGRAM_BOT_KEY** environment variable. You must set the bot's webhook to the */endpoints/telegram* endpoint in order to receive messages.

<a name="examples"></a>
## 12. Examples

I have provided some config files unnder *examples*. Run ```chatto``` with the ```--path``` of your desired example to test them out.

1. [**Mood Bot**](/examples/01_moodbot) - Greet the bot to start the conversation.
2. [**Engineering Flowchart**](/examples/02_repair) - Tell the bot you want to repair something.
3. [**Pokemon**](/examples/03_pokemon) - Search for Pokémon by name or number.
4. [**Trivia Quiz**](/examples/04_trivia) - Type *start* to take a quick trivia quiz.
