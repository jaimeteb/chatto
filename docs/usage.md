# Usage

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

!!! note
    The default log level is **INFO**. You can set it to **DEBUG** with the environment variable `DEBUG` set to `true`.

## CLI

You can use the Chatto CLI tool by downloading the `chatto-cli` binary. The CLI makes it easy to test your bot interactions.

```bash
chatto-cli -url 'http://localhost' -port 4770
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

## Chatto Init

Download the `chatto-init` binary to generate a sample Chatto project on the directory of your choice. Just run:

```bash
chatto-init -path my-chatto/
```

A Chatto project will be initialized at `my-chatto`.

## Import

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

## Deployment
 
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

!!! important
    The [extensions](/extension) server has to be executed according to its language.
    For this `docker-compose.yml` file, you'd have to build the Go extension first:

    ```bash
    go build -o data/ext/ext data/ext/ext.go
    ```

!!! note
    The [extensions](/extension) server has to be running before Chatto initializes.

### Kubernetes

Under the `deploy/kubernetes` directory you can find an example deployment:

| Kind       | Name                    | Description                                                   |
|------------|-------------------------|---------------------------------------------------------------|
| Secret     | `chatto-config-secrets` | Contains the tokens that Chatto will use for authorization    |
| ConfigMap  | `chatto-config-envs`    | Contains the environment variables for the **bot.yml** file   |
| ConfigMap  | `chatto-config-files`   | Contains the **clf.yml** and **fsm.yml** file                 |
| Deployment | `chatto`                | Chatto deployment based on the `jaimeteb/chatto` Docker image |
| Service    | `chatto-service`        | Service for the `chatto` deployment                           |
| Ingress    | `chatto-ingress`        | Ingress for the `chatto-service` service                      |

Run the following command to deploy on Kubernetes:

```bash
kubectl apply -f ./deploy/kubernetes/
```
