package logging

import (
	"github.com/rancher/logging-controller/manager"
	"github.com/rancher/types/apis/logging.cattle.io/v3"
)

func Register(ma *manager.Manager) {
	lifecycle := &LoggingLifecycle{
		Manager: ma,
	}

	loggingClient := ma.LogCtx.Logging.Loggings("")
	loggingClient.AddLifecycle("logging-controller", lifecycle)
}

type LoggingLifecycle struct {
	Manager *manager.Manager
}

func (c *LoggingLifecycle) Create(obj *v3.Logging) (*v3.Logging, error) {
	return nil, nil
}

func (c *LoggingLifecycle) Remove(obj *v3.Logging) (*v3.Logging, error) {
	// c.Manager.Stop(obj)
	return nil, nil
}

func (c *LoggingLifecycle) Updated(obj *v3.Logging) (*v3.Logging, error) {
	err := c.Manager.ReloadConfig("project")
	return obj, err
}
