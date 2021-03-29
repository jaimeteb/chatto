# Getting Started

With Chatto you can create simple and fast chatbots, starting with a few YAML configuration files.

If you're using Go, you can install it via `go get`. There are [Docker images](https://hub.docker.com/repository/docker/jaimeteb/chatto) and [binaries](https://github.com/jaimeteb/chatto/releases) available as well.

## Installation

### Via go get

To download the latest version:

```
go get -u github.com/jaimeteb/chatto
```

To install a specific version, initialize a **Go Module** and run:

```
go get -u github.com/jaimeteb/chatto/...@version
```

### Via Docker

```
docker pull jaimeteb/chatto:latest
```

You can specify a version as well:

```
docker pull jaimeteb/chatto:version
```