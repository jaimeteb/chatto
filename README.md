[![Build Status](https://travis-ci.com/jaimeteb/chatto.svg?branch=master)](https://travis-ci.com/jaimeteb/chatto)
[![codecov](https://codecov.io/gh/jaimeteb/chatto/branch/master/graph/badge.svg)](https://codecov.io/gh/jaimeteb/chatto)
[![Go Report Card](https://goreportcard.com/badge/github.com/jaimeteb/chatto)](https://goreportcard.com/report/github.com/jaimeteb/chatto)
[![GoDoc](https://godoc.org/github.com/jaimeteb/chatto?status.svg)](https://godoc.org/github.com/jaimeteb/chatto)
---
# chatto

<p align="center">
<img src="https://user-images.githubusercontent.com/17936011/89082867-e3c0d300-d354-11ea-9def-008c403a4497.png" alt="botto" width="150"/>
</p>

Simple chatbot framework written in Go, with configurations in YAML. The aim of this project is to create very simple text-based chatbots using a few configuration files. 

The inspiration for this project came from [Flottbot](https://github.com/target/flottbot), and my experience using [Rasa](https://github.com/RasaHQ/rasa).

## Installation

Run ```go get -u github.com/jaimeteb/chatto```.

## How does it work?

Chatto combines the consistency of a finite-state-machine with the flexibility of machine learning. It has two main components: the classifier and the finite-stete-machine

### 1. Classifier

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

### 2. Finite State Machine

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
    message: "Turning off."

defaults:
  unknown: "Can't do that."
  unsure: "???"
```

Under **functions** you can list the transitions available for the FSM. The object **transition** describes the states of the transition (**from** one state **into** another) if **command** is executed; **message** is the message to send to the user.

The special state **any** can help you to go from any state into another, if the command is executed. You don't have to declare the **any** state in the states list.

## 3. Extensions

The extensions in chatto are pieces of code that can be executed instead of messages. The extensions names must begin by **"ext_"** and they must be placed in the **ext/ext.go** file. The format for a chatto extension is as follows:

```go
package main

import (
	"github.com/jaimeteb/chatto/fsm"
)

func greetFunc(m *fsm.FSM) interface{} {
	return "Hello Universe"
}

// Ext is exported
var Ext = fsm.FuncMap{
	"ext_any": greetFunc,
}

func main() {}
```

You must export a variable called **Ext** of type **fsm.FuncMap**, and map the extension names to their respective functions. In this example, **ext_any** simply returns "Hello Universe". Note that, the functions must have the ```func(*FSM) interface{}``` signature.

## 4. Slots

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

## HTTP Endpoint

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

## CLI

Alternatively, run chatto on a command line interface using the ```--cli``` flag.

## Examples

I have provided some config files unnder *examples*. Run ```chatto``` with the ```--path``` of your desired example to test them out.

1. [**Mood Bot**](/examples/01_moodbot) - Greet the bot to start the conversation.
2. [**Engineering Flowchart**](/examples/02_repair) - Tell the bot you want to repair something.
3. [**Echo**](/examples/03_echo) - The bot will ask for your name and age (beta).
4. [**Trivia Quiz**](/examples/04_trivia) - Type *start* to take a quick trivia quiz.
