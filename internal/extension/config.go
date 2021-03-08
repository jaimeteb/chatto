package extension

import (
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

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
func New(extCfg []Config) (Map, error) {
	extensionMap := make(Map)

	for n := range extCfg {
		ext := extCfg[n]

		switch ext.Type {
		case "RPC":
			client, err := dialRPC(ext.Host, ext.Port)
			if err != nil {
				log.Errorf("unable to get rpc extensions for '%s:%d': %s", ext.Host, ext.Port, err)
				continue
			}

			rpcExtension := &RPC{client}

			extensionNames, err := rpcExtension.GetAllExtensions()
			if err != nil {
				log.Errorf("unable to get rpc extensions for '%s:%d': %s", ext.Host, ext.Port, err)
				continue
			}

			for n := range extensionNames {
				addErr := extensionMap.Add(extensionNames[n], rpcExtension)
				if addErr != nil {
					return nil, addErr
				}
			}
		case "REST":
			retryClient := retryablehttp.NewClient()
			retryClient.Logger = nil

			restExtension := &REST{URL: ext.URL, http: retryClient, token: ext.Token}

			extensionNames, err := restExtension.GetAllExtensions()
			if err != nil {
				log.Errorf("unable to get rest extensions for '%s:%d': %s", ext.URL, ext.Port, err)
				continue
			}

			for n := range extensionNames {
				addErr := extensionMap.Add(extensionNames[n], restExtension)
				if addErr != nil {
					return nil, addErr
				}
			}
		default:
			return nil, fmt.Errorf("invalid extension type: %s", ext.Type)
		}
	}

	log.Info("Loaded extensions:")
	for name, _ := range extensionMap {
		log.Infof("%s", name)
	}

	if len(extensionMap) == 0 {
		log.Info("Using bot without extensions.")
	}

	return extensionMap, nil
}
