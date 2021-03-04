[![Documentation](https://img.shields.io/static/v1?label=&message=Documentation&color=red)](https://chatto.jaimeteb.com)
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

<p align="center">
<img src="https://media.giphy.com/media/DFIxYClozxyMg9wnil/giphy.gif" alt="demo" width="480"/>
</p>

## Contents

* [Installation](#install)
* [Documentation](#docs)
* [Your first bot](#yourfirstbot)
    * [The **clf.yml** file](#yourfirstbotclf)
    * [The **fsm.yml** file](#yourfirstbotfsm)
    * [Run your bot](#yourfirstbotrun)
    * [Interact with your bot](#yourfirstbotinteract)
* [Usage](#usage)  
    * [CLI](#usagecli)
    * [Docker Compose](#usagecompose)
    * [Import](#usageimport)
* [Examples](#examples)  

<a name="install"></a>
## Installation

```
go get -u github.com/jaimeteb/chatto/cmd/chatto
```

Via Docker:

```
docker pull jaimeteb/chatto:latest
```

<a name="docs"></a>
## Documentation

See the [**Documentation**](https://chatto.jaimeteb.com) for **examples**, **configuration guides** and **reference**.

<p align="center">
<img src="https://i.imgur.com/RkgEfX2.jpg" href="https://chatto.jaimeteb.com" alt="docs"/>
</p>

<a name="yourfirstbot"></a>
## Your first bot

Chatto combines the consistency of a finite-state-machine with the flexibility of machine learning. It has three main components: the classifier, the finite-state-machine and the extensions.

A very basic directory structure for Chatto would be the following:

```
.
└──data
   ├── clf.yml
   └── fsm.yml
```

Start by creating the `data` directory as well as the YAML files.

```console
mkdir data
touch data/clf.yml data/fsm.yml
```

<a name="yourfirstbotclf"></a>
### The **clf.yml** file

The **clf.yml** file defines how the user messages will be classified into *commands* (intents). Start with this very simple configuration:

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

<a name="yourfirstbotfsm"></a>
### The **fsm.yml** file

The **fsm.yml** file defines the transitions between states, the commands that make these transitions, and the answers to be sent in them. Start with this file contents:

```yaml
states:
  - "off"
  - "on"

commands:
  - "turn_on"
  - "turn_off"

transitions:
  - from:
      - "off"
    into: "on"
    command: "turn_on"
    answers:
      - text: "Turning on."

  - from:
      - "on"
    into: "off"
    command: "turn_off"
    answers:
      - text: "Turning off."
      - text: "❌"

defaults:
  unknown: "Can't do that."
```

<a name="yourfirstbotrun"></a>
### Run your first bot

To start your bot, run:

```bash
chatto -path data/
```

If you're using Docker, run:

```bash
docker run \
    -it \
    -e CHATTO_DATA=./data \
    -v $PWD/data:/data \
    jaimeteb/chatto:latest \
    chatto -path data
```

<a name="yourfirstbotinteract"></a>
### Interact with your first bot

To interact with your bot, run:

```
chatto-cli
```

That's it! Now you can say *turn on* or *on* to go into the **on** state, and *turn off* or *off* to go back into **off**. However, you cannot go from **on** into **on**, or from **off** into **off** either.

Here is a diagram for this simple Finite State Machine:

![ON/OFF Finite State Machine](https://uploads.gamedev.net/monthly_06_2013/ccs-209764-0-84996300-1370053229.jpg)


<a name="usage"></a>
## Usage

> You can integrate your bot with [**Telegram, Twilio, Slack**](https://chatto.jaimeteb.com/channels/) and [**anything you like**](https://chatto.jaimeteb.com/endpoints/)

Run `chatto` in the directory where your YAML files are located, or specify a path to them with the `-path` flag:

```bash
chatto -path ./your/data
```

To run on Docker, use:

```bash
docker run \
  -p 4770:4770 \
  -e CHATTO_DATA=./your/data \
  -v $PWD/your/data:/data \
  jaimeteb/chatto
```

<a name="usagecli"></a>
### CLI

You can use the Chatto CLI tool by downloading the `chatto-cli` binary. The CLI makes it easy to test your bot interactions.

```bash
chatto-cli -url 'http://mybot.com' -port 4770
```

<a name="usagecompose"></a>
### Docker Compose

You can use Chatto on Docker Compose as well. A `docker-compose.yml` would look like this:

```yaml
version: "3"

services:
  chatto:
    image: jaimeteb/chatto:${CHATTO_VERSION}
    env_file: .env
    ports:
      - "4770:4770"
    volumes:
      - ${CHATTO_DATA}:/data
    depends_on:
      - ext
      - redis

  ext:
    image: odise/busybox-curl # Busy box with certificates
    command: ext/ext
    expose:
      - 8770
    volumes:
      - ${CHATTO_DATA}/ext:/ext

  redis:
    image: bitnami/redis:6.0
    environment:
      - REDIS_PASSWORD=${STORE_PASSWORD}
    expose:
      - 6379
```

This requires a `.env` file to contain the necessary environment variables:

```
# Chatto configuration
CHATTO_VERSION=latest
CHATTO_DATA=./your/data

# Extension configuration
EXTENSIONS_URL=http://ext:8770

# Redis
STORE_HOST=redis
STORE_PASSWORD=pass

# Logs
DEBUG=true
```

The directory structure with all the files would look like this:

```
.
├── data
│   ├── ext
│   │   ├── ext
│   │   └── ext.go
│   ├── bot.yml
│   ├── chn.yml
│   ├── clf.yml
|   └── fsm.yml
├── docker-compose.yml
└── .env
```

Finally, run:

```bash
docker-compose up -d redis ext
docker-compose up -d chatto
```

> The [extensions](/extensions) server has to be executed according to its language.<br><br>For this `docker-compose.yml` file, you'd have to build the Go extension first:<br><br>```go build -o data/ext/ext data/ext/ext.go```

> The [extensions](/extensions) server has to be running before Chatto initializes.

<a name="usageimport"></a>
### Import

An importable bot server and client package is provided to allow embedding into your own application.

To embed the server:

```go
package main

import (
	"flag"

	"github.com/jaimeteb/chatto/bot"
)

func main() {
	port := flag.Int("port", 4770, "Specify port to use.")
	path := flag.String("path", ".", "Path to YAML files.")
	flag.Parse()

	server := bot.NewServer(*path, *port)

	server.Run()
}
```

To embed the client:

```go
package myservice

import (
	"log"

	"github.com/jaimeteb/chatto/bot"
)

type MyService struct {
	chatto bot.Client
}

func NewMyService(url string, port int) *MyService {
	return &MyService{chatto: bot.NewClient(url, port)}
}

func (s *MyService) Submit(question *query.Question) error {
	answers, err := s.chatto.Submit(question)
	if err != nil {
		return err
	}

	// Print answers to stdout
	for _, answer := range answers {
		fmt.Println(answer.Text)
	}

	return nil
}
```

<a name="examples"></a>
## Examples

I have provided some config files under *examples*. Clone the repository and run `chatto` with the `-path` of your desired example to test them out (for the ones that use extensions, run their respective extensions first).

More about these examples in the [**Documentation**](https://chatto.jaimeteb.com/examples/moodbot)

1. [**Mood Bot**](/examples/01_moodbot) - A chatto version of [Rasa's Mood Bot](https://github.com/RasaHQ/rasa/tree/master/examples/moodbot) Greet the bot to start the conversation.
3. [**Pokemon Search**](/examples/03_pokemon) - Search for Pokémon by name or number.
2. [**Miscellaneous Bot**](/examples/02_misc) - Weather forecast, random jokes and quotes, and more!
4. [**Trivia Quiz**](/examples/04_trivia) - Type *start* to take a quick trivia quiz.
