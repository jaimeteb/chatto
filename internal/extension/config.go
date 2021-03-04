package extension

import (
	"fmt"
	"net"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

// Config options for an extension function
type Config struct {
	Type  string `mapstructure:"type"`
	Host  string `mapstructure:"host"`
	Port  int    `mapstructure:"port"`
	URL   string `mapstructure:"url"`
	Token string `mapstructure:"token"`
}

// New loads the extension configuration and connects to the server
func New(extCfg Config) (Extension, error) {
	var extension Extension

	switch extCfg.Type {
	case "RPC":
		maxTries := 3
		backOff := 1 * time.Second
		var tries int
		var client *rpc.Client
		var err error
		for {
			client, err = rpc.Dial("tcp", fmt.Sprintf("%s:%d", extCfg.Host, extCfg.Port))
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

		rpcExtension := &RPC{client}

		allFuncs, err := rpcExtension.GetAllExtensions()
		if err != nil {
			return nil, err
		}

		log.Info("Loaded extensions (RPC):")
		for i, fun := range allFuncs {
			log.Infof("%2d %v", i, fun)
		}

		extension = rpcExtension
	case "REST":
		retryClient := retryablehttp.NewClient()
		retryClient.Logger = nil

		restExtention := &REST{URL: extCfg.URL, http: retryClient, token: extCfg.Token}

		allFuncs, err := restExtention.GetAllExtensions()
		if err != nil {
			return nil, err
		}

		log.Info("Loaded extensions (REST):")
		for i, fun := range allFuncs {
			log.Infof("%2d %v", i, fun)
		}

		extension = restExtention
	}

	if extension == nil {
		log.Info("Using bot without extensions.")
	}

	return extension, nil
}
