package controller

import (
	"github.com/rancher/logging-controller/controller/projectlogging"
	"github.com/rancher/logging-controller/controller/logging"
	"github.com/rancher/logging-controller/manager"
)

func Register(ma *manager.Manager) {
	logging.Register(ma)
	projectlogging.Register(ma)
}
