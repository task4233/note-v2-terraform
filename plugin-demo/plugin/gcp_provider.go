package main

import (
	"fmt"
	"os"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	provider "github.com/task4233/go-plugin/common"
)

type GCPProvider struct {
	logger hclog.Logger
}

func (p *GCPProvider) Create(scheme *provider.Scheme) (string, error) {
	p.logger.Debug("Creating GCP Resources ...")
	val, ok := scheme.Default.(string)
	if !ok {
		return "", fmt.Errorf("failed conversion: %v", scheme.Default)
	}
	p.logger.Debug("message: " + val)

	return "Created", nil
}

func (p *GCPProvider) Update(scheme *provider.Scheme) (string, error) {
	p.logger.Debug("Updating GCP Resources ...")
	val, ok := scheme.Default.(string)
	if !ok {
		return "", fmt.Errorf("failed conversion: %v", scheme.Default)
	}
	p.logger.Debug("message: " + val)

	return "Updated", nil
}

func (p *GCPProvider) Delete(scheme *provider.Scheme) (string, error) {
	p.logger.Debug("Deleting GCP Resources ...")
	val, ok := scheme.Default.(string)
	if !ok {
		return "", fmt.Errorf("failed conversion: %v", scheme.Default)
	}
	p.logger.Debug("message: " + val)

	return "Deleted", nil
}

func (p *GCPProvider) Get(id string, resp *provider.Scheme) error {
	p.logger.Debug("Getting GCP Resources ...")
	*resp = provider.Scheme{
		Default: "GCP vm",
	}
	return nil
}

var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PROVIDER_PLUGIN",
	MagicCookieValue: "gcp",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	gcpProvider := &GCPProvider{
		logger: logger,
	}

	var pluginMap = map[string]plugin.Plugin{
		"gcp": &provider.ProviderPlugin{Impl: gcpProvider},
	}
	logger.Debug("Now GCP Provider Serving")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
