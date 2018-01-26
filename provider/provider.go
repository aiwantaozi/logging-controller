package provider

import (
	"fmt"

	"github.com/rancher/types/config"
	"github.com/urfave/cli"
)

type LogProvider interface {
	Init(*cli.Context, *config.LoggingContext)
	Start()
	Stop() error
	Reload() error
	GetName() string
	GetConfig(name string) LogProviderConfig
}

type LogProviderConfig interface {
	Update() error
}

var (
	providers map[string]LogProvider
)

func GetProvider(name string, c *cli.Context, lc *config.LoggingContext) LogProvider {
	if provider, ok := providers[name]; ok {
		provider.Init(c, lc)
		return provider
	}
	return providers["fluentd"]
}

func RegisterProvider(name string, provider LogProvider) error {
	if providers == nil {
		providers = make(map[string]LogProvider)
	}
	if _, exists := providers[name]; exists {
		return fmt.Errorf("provider already registered")
	}
	providers[name] = provider
	return nil
}
