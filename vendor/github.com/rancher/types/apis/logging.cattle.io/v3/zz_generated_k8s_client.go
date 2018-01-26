package v3

import (
	"context"
	"sync"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/controller"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type Interface interface {
	RESTClient() rest.Interface
	controller.Starter

	LoggingsGetter
	ProjectLoggingsGetter
}

type Client struct {
	sync.Mutex
	restClient rest.Interface
	starters   []controller.Starter

	loggingControllers        map[string]LoggingController
	projectLoggingControllers map[string]ProjectLoggingController
}

func NewForConfig(config rest.Config) (Interface, error) {
	if config.NegotiatedSerializer == nil {
		configConfig := dynamic.ContentConfig()
		config.NegotiatedSerializer = configConfig.NegotiatedSerializer
	}

	restClient, err := rest.UnversionedRESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &Client{
		restClient: restClient,

		loggingControllers:        map[string]LoggingController{},
		projectLoggingControllers: map[string]ProjectLoggingController{},
	}, nil
}

func (c *Client) RESTClient() rest.Interface {
	return c.restClient
}

func (c *Client) Sync(ctx context.Context) error {
	return controller.Sync(ctx, c.starters...)
}

func (c *Client) Start(ctx context.Context, threadiness int) error {
	return controller.Start(ctx, threadiness, c.starters...)
}

type LoggingsGetter interface {
	Loggings(namespace string) LoggingInterface
}

func (c *Client) Loggings(namespace string) LoggingInterface {
	objectClient := clientbase.NewObjectClient(namespace, c.restClient, &LoggingResource, LoggingGroupVersionKind, loggingFactory{})
	return &loggingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}

type ProjectLoggingsGetter interface {
	ProjectLoggings(namespace string) ProjectLoggingInterface
}

func (c *Client) ProjectLoggings(namespace string) ProjectLoggingInterface {
	objectClient := clientbase.NewObjectClient(namespace, c.restClient, &ProjectLoggingResource, ProjectLoggingGroupVersionKind, projectLoggingFactory{})
	return &projectLoggingClient{
		ns:           namespace,
		client:       c,
		objectClient: objectClient,
	}
}
