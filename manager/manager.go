package manager

import (
	"github.com/rancher/logging-controller/provider"
	"github.com/rancher/types/config"
	"github.com/urfave/cli"
)

type Manager struct {
	LogCtx   *config.LoggingContext
	provider provider.LogProvider
}

type Secret struct {
	TargetType string                 `json:"type"`
	Label      string                 `json:"label"`
	Data       map[string]interface{} `json:"data"`
}

func New(c *cli.Context, logCtx *config.LoggingContext) *Manager {
	provider := provider.GetProvider("fluentd", c, logCtx)
	return &Manager{
		LogCtx:   logCtx,
		provider: provider,
	}

}

func (m *Manager) Start() error {
	conf := m.provider.GetConfig("cluster")
	err := conf.Update()
	// .Update()
	if err != nil {
		return err
	}
	err = m.provider.GetConfig("project").Update()
	if err != nil {
		return err
	}
	m.provider.Start()
	return nil
}

func (m *Manager) ReloadConfig(name string) error {
	err := m.provider.GetConfig(name).Update()
	if err != nil {
		return err
	}
	err = m.provider.Reload()
	if err != nil {
		return err
	}
	return nil
}
