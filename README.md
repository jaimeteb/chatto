[![Documentation](https://img.shields.io/static/v1?label=&message=Documentation&color=red)](https://chatto.jaimeteb.com)
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
* [Usage](#usage)  
    * [CLI](#usagecli)
    * [Docker Compose](#usagecompose)
* [Examples](#examples)  

<a name="install"></a>
## Installation

```
go get -u github.com/jaimeteb/chatto
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

The **fsm.yml** file defines the transitions between states, the commands that make these transitions, and the messages to be sent in them. Start with this file contents:

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

<a name="yourfirstbotrun"></a>
### Run your first bot

To run your bot in a CLI, simply run:

```bash
chatto -cli -path data/
```

Or if you're using Docker, run:

```bash
docker run \
    -it \
    -e CHATTO_DATA=./data \
    -v $PWD/data:/chatto/data \
    jaimeteb/chatto:latest \
    chatto -cli -path data
```

That's it! Now you can say *turn on* or *on* to go into the **on** state, and *turn off* or *off* to go back into **off**. However, you cannot go from **on** into **on**, or from **off** into **off** either.

Here is a diagram for this simple Finite State Machine:

![ON/OFF Finite State Machine](https://uploads.gamedev.net/monthly_06_2013/ccs-209764-0-84996300-1370053229.jpg)


<a name="usage"></a>
## Usage

> You can integrate yout bot with [**Telegram and Twilio**](https://chatto.jaimeteb.com/botconfiguration/) and [**anything you like**](https://chatto.jaimeteb.com/endpoints/)

Run `chatto` in the directory where your YAML files are located, or specify a path to them with the `-path` flag:

```bash
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

> You can set a log level with the environment variable `LOG_LEVEL`.

<a name="usagecli"></a>
### CLI

You can use Chatto in a CLI mode by adding the `-cli` flag.

```bash
chatto -cli -path ./your/data
```

On Docker:

```bash
docker run \
    -it \
    -e CHATTO_DATA=./your/data \
    -v $PWD./your/data:/chatto/data \
    jaimeteb/chatto:latest \
    chatto -cli -path data
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
      - ${CHATTO_DATA}:/chatto/data
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
LOG_LEVEL=DEBUG
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


<a name="examples"></a>
## Examples

I have provided some config files under *examples*. Clone the repository and run `chatto` with the `-path` of your desired example to test them out (for the ones that use extensions, run their respective extensions first).

More about these examples in the [**Documentation**](https://chatto.jaimeteb.com/examples/moodbot)

1. [**Mood Bot**](/examples/01_moodbot) - A chatto version of [Rasa's Mood Bot](https://github.com/RasaHQ/rasa/tree/master/examples/moodbot) Greet the bot to start the conversation.
2. [**Engineering Flowchart**](/examples/02_repair) - Tell the bot you want to repair something.
3. [**Pokemon**](/examples/03_pokemon) - Search for Pokémon by name or number.
4. [**Trivia Quiz**](/examples/04_trivia) - Type *start* to take a quick trivia quiz.
