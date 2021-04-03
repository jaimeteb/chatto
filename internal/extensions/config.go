package extensions

import (
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

// ConfigMap of extension server names to their configs
type ConfigMap map[string]Config

// Config contains all the require parameters
// to communicate with an extension
type Config struct {
	Type  string `mapstructure:"type"`
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

func dialRPC(host string, port int) (*rpc.Client, error) {
	maxTries := 3
	backOff := 1 * time.Second
	var tries int
	var client *rpc.Client
	var err error
	for {
		client, err = rpc.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			switch err.(type) {
			case *net.OpError:
				if tries > maxTries {
					return nil, err
				}

				tries++

				time.Sleep(backOff)

				continue
			default:
				return nil, err
			}
		}

		break
	}

	return client, nil
}

// New loads the extension configuration and connects to the server
func New(extCfg ConfigMap) (*WebSocketServer, ServerMap, error) {
	extensionMap := make(ServerMap)

	webSocket := NewWebSocket()

	for server, config := range extCfg {
		switch config.Type {
		case "RPC":
			client, err := dialRPC(config.Host, config.Port)
			if err != nil {
				log.Errorf("unable to get rpc extensions for '%s:%d': %s", config.Host, config.Port, err)
				continue
			}

			rpcExtension := &RPC{client}

			addErr := extensionMap.Add(server, rpcExtension)
			if addErr != nil {
				return nil, nil, addErr
			}
		case "REST":
			retryClient := retryablehttp.NewClient()
			retryClient.Logger = nil

			restExtension := &REST{URL: config.URL, http: retryClient, token: config.Token}

			addErr := extensionMap.Add(server, restExtension)
			if addErr != nil {
				return nil, nil, addErr
			}
		case "WEBSOCKET":
			addErr := extensionMap.Add(server, webSocket)
			if addErr != nil {
				return nil, nil, addErr
			}
		default:
			return nil, nil, fmt.Errorf("invalid extension type: %s", config.Type)
		}
	}

	log.Info("Loaded extensions:")
	for name := range extensionMap {
		log.Infof("%s", name)
	}

	if len(extensionMap) == 0 {
		log.Info("Using bot without extensions.")
	}

	return webSocket, extensionMap, nil
}
