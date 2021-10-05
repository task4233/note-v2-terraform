package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	provider "github.com/task4233/go-plugin/common"
)

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PROVIDER_PLUGIN",
	MagicCookieValue: "gcp",
}

var pluginMap = map[string]plugin.Plugin{
	"gcp": &provider.ProviderPlugin{},
}

func run() int {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "client",
		Level:  hclog.Debug,
		Output: os.Stdout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  handshakeConfig,
		Plugins:          pluginMap,
		VersionedPlugins: map[int]plugin.PluginSet{},
		Cmd:              exec.CommandContext(ctx, "./plugin/plugin"),
		Logger:           logger,
	})
	defer client.Kill()

	rpcClient, err := client.Client()
	if err != nil {
		log.Printf("failed client.Client: %s", err.Error())
		return 1
	}

	raw, err := rpcClient.Dispense("gcp")
	if err != nil {
		log.Printf("failed client.Dispence: %s", err.Error())
		return 1
	}

	gcpProvider, ok := raw.(provider.Provider)
	if !ok {
		log.Printf("failed conversion: %v", raw)
		return 1
	}

	message := "Hello~"
	scheme := &provider.Scheme{
		Default: message,
	}

	// Create
	{
		mes, err := gcpProvider.Create(scheme)
		if err != nil {
			log.Printf("failed Create: %s", err.Error())
			return 1
		}
		log.Println(mes)
	}

	// Update
	{
		scheme.Default = "Hello2"
		mes, err := gcpProvider.Update(scheme)
		if err != nil {
			log.Printf("failed Update: %s", err.Error())
			return 1
		}
		log.Println(mes)
	}

	// Get
	{
		// might be wrong?
		id := "Hello2"
		err := gcpProvider.Get(id, scheme)
		if err != nil {
			log.Printf("failed Get: %s", err.Error())
			return 1
		}
		log.Println(scheme)
	}

	// Delete
	{
		mes, err := gcpProvider.Delete(scheme)
		if err != nil {
			log.Printf("failed Delete: %s", err.Error())
			return 1
		}
		log.Println(mes)
	}

	return 0
}

func main() {
	os.Exit(run())
}
