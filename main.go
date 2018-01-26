package main

import (
	"os"

	"github.com/rancher/logging-controller/controller"
	"github.com/rancher/logging-controller/manager"
	"github.com/rancher/types/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"k8s.io/client-go/tools/clientcmd"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "logging-controller"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug",
		},
		cli.StringFlag{
			Name:   "config",
			Usage:  "Kube config for accessing kubernetes cluster",
			EnvVar: "KUBECONFIG",
		},
		cli.BoolFlag{
			Name:  "fluentd-dry-run",
			Usage: "generate the config file, but not run fluentd",
		},
		cli.StringFlag{
			Name:  "fluentd-config-dir",
			Usage: "Fluentd config directory",
			Value: "/fluentd/etc/config",
		},
	}

	app.Action = func(c *cli.Context) error {
		logrus.Info("I'm a turkey")
		return nil
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", c.String("config"))
	if err != nil {
		return err
	}

	logCtx, err := config.NewLoggingContext(*kubeConfig)
	if err != nil {
		return err
	}
	ma := manager.New(c, logCtx)
	ma.Start()
	controller.Register(ma)
	return nil
}
