package projectlogging

import (
	"github.com/rancher/logging-controller/manager"
	"github.com/rancher/types/apis/logging.cattle.io/v3"
)

func Register(ma *manager.Manager) {
	lifecycle := &ProjectLoggingLifecycle{
		Manager: ma,
	}

	projectloggingClient := ma.LogCtx.Logging.ProjectLoggings("")
	projectloggingClient.AddLifecycle("project-logging-controller", lifecycle)
}

type ProjectLoggingLifecycle struct {
	Manager *manager.Manager
}

func (c *ProjectLoggingLifecycle) Create(obj *v3.ProjectLogging) (*v3.ProjectLogging, error) {
	return nil, nil
}

func (c *ProjectLoggingLifecycle) Remove(obj *v3.ProjectLogging) (*v3.ProjectLogging, error) {
	// c.Manager.Stop(obj)
	return nil, nil
}

func (c *ProjectLoggingLifecycle) Updated(obj *v3.ProjectLogging) (*v3.ProjectLogging, error) {
	err := c.Manager.ReloadConfig("project")
	return obj, err
}
