package v3

import (
	"context"

	"github.com/rancher/norman/clientbase"
	"github.com/rancher/norman/controller"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

var (
	LoggingGroupVersionKind = schema.GroupVersionKind{
		Version: Version,
		Group:   GroupName,
		Kind:    "Logging",
	}
	LoggingResource = metav1.APIResource{
		Name:         "loggings",
		SingularName: "logging",
		Namespaced:   false,
		Kind:         LoggingGroupVersionKind.Kind,
	}
)

type LoggingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Logging
}

type LoggingHandlerFunc func(key string, obj *Logging) error

type LoggingLister interface {
	List(namespace string, selector labels.Selector) (ret []*Logging, err error)
	Get(namespace, name string) (*Logging, error)
}

type LoggingController interface {
	Informer() cache.SharedIndexInformer
	Lister() LoggingLister
	AddHandler(name string, handler LoggingHandlerFunc)
	AddClusterScopedHandler(name, clusterName string, handler LoggingHandlerFunc)
	Enqueue(namespace, name string)
	Sync(ctx context.Context) error
	Start(ctx context.Context, threadiness int) error
}

type LoggingInterface interface {
	ObjectClient() *clientbase.ObjectClient
	Create(*Logging) (*Logging, error)
	GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Logging, error)
	Get(name string, opts metav1.GetOptions) (*Logging, error)
	Update(*Logging) (*Logging, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error
	List(opts metav1.ListOptions) (*LoggingList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Controller() LoggingController
	AddHandler(name string, sync LoggingHandlerFunc)
	AddLifecycle(name string, lifecycle LoggingLifecycle)
	AddClusterScopedHandler(name, clusterName string, sync LoggingHandlerFunc)
	AddClusterScopedLifecycle(name, clusterName string, lifecycle LoggingLifecycle)
}

type loggingLister struct {
	controller *loggingController
}

func (l *loggingLister) List(namespace string, selector labels.Selector) (ret []*Logging, err error) {
	err = cache.ListAllByNamespace(l.controller.Informer().GetIndexer(), namespace, selector, func(obj interface{}) {
		ret = append(ret, obj.(*Logging))
	})
	return
}

func (l *loggingLister) Get(namespace, name string) (*Logging, error) {
	var key string
	if namespace != "" {
		key = namespace + "/" + name
	} else {
		key = name
	}
	obj, exists, err := l.controller.Informer().GetIndexer().GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(schema.GroupResource{
			Group:    LoggingGroupVersionKind.Group,
			Resource: "logging",
		}, name)
	}
	return obj.(*Logging), nil
}

type loggingController struct {
	controller.GenericController
}

func (c *loggingController) Lister() LoggingLister {
	return &loggingLister{
		controller: c,
	}
}

func (c *loggingController) AddHandler(name string, handler LoggingHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}
		return handler(key, obj.(*Logging))
	})
}

func (c *loggingController) AddClusterScopedHandler(name, cluster string, handler LoggingHandlerFunc) {
	c.GenericController.AddHandler(name, func(key string) error {
		obj, exists, err := c.Informer().GetStore().GetByKey(key)
		if err != nil {
			return err
		}
		if !exists {
			return handler(key, nil)
		}

		if !controller.ObjectInCluster(cluster, obj) {
			return nil
		}

		return handler(key, obj.(*Logging))
	})
}

type loggingFactory struct {
}

func (c loggingFactory) Object() runtime.Object {
	return &Logging{}
}

func (c loggingFactory) List() runtime.Object {
	return &LoggingList{}
}

func (s *loggingClient) Controller() LoggingController {
	s.client.Lock()
	defer s.client.Unlock()

	c, ok := s.client.loggingControllers[s.ns]
	if ok {
		return c
	}

	genericController := controller.NewGenericController(LoggingGroupVersionKind.Kind+"Controller",
		s.objectClient)

	c = &loggingController{
		GenericController: genericController,
	}

	s.client.loggingControllers[s.ns] = c
	s.client.starters = append(s.client.starters, c)

	return c
}

type loggingClient struct {
	client       *Client
	ns           string
	objectClient *clientbase.ObjectClient
	controller   LoggingController
}

func (s *loggingClient) ObjectClient() *clientbase.ObjectClient {
	return s.objectClient
}

func (s *loggingClient) Create(o *Logging) (*Logging, error) {
	obj, err := s.objectClient.Create(o)
	return obj.(*Logging), err
}

func (s *loggingClient) Get(name string, opts metav1.GetOptions) (*Logging, error) {
	obj, err := s.objectClient.Get(name, opts)
	return obj.(*Logging), err
}

func (s *loggingClient) GetNamespaced(namespace, name string, opts metav1.GetOptions) (*Logging, error) {
	obj, err := s.objectClient.GetNamespaced(namespace, name, opts)
	return obj.(*Logging), err
}

func (s *loggingClient) Update(o *Logging) (*Logging, error) {
	obj, err := s.objectClient.Update(o.Name, o)
	return obj.(*Logging), err
}

func (s *loggingClient) Delete(name string, options *metav1.DeleteOptions) error {
	return s.objectClient.Delete(name, options)
}

func (s *loggingClient) DeleteNamespaced(namespace, name string, options *metav1.DeleteOptions) error {
	return s.objectClient.DeleteNamespaced(namespace, name, options)
}

func (s *loggingClient) List(opts metav1.ListOptions) (*LoggingList, error) {
	obj, err := s.objectClient.List(opts)
	return obj.(*LoggingList), err
}

func (s *loggingClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return s.objectClient.Watch(opts)
}

// Patch applies the patch and returns the patched deployment.
func (s *loggingClient) Patch(o *Logging, data []byte, subresources ...string) (*Logging, error) {
	obj, err := s.objectClient.Patch(o.Name, o, data, subresources...)
	return obj.(*Logging), err
}

func (s *loggingClient) DeleteCollection(deleteOpts *metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	return s.objectClient.DeleteCollection(deleteOpts, listOpts)
}

func (s *loggingClient) AddHandler(name string, sync LoggingHandlerFunc) {
	s.Controller().AddHandler(name, sync)
}

func (s *loggingClient) AddLifecycle(name string, lifecycle LoggingLifecycle) {
	sync := NewLoggingLifecycleAdapter(name, false, s, lifecycle)
	s.AddHandler(name, sync)
}

func (s *loggingClient) AddClusterScopedHandler(name, clusterName string, sync LoggingHandlerFunc) {
	s.Controller().AddClusterScopedHandler(name, clusterName, sync)
}

func (s *loggingClient) AddClusterScopedLifecycle(name, clusterName string, lifecycle LoggingLifecycle) {
	sync := NewLoggingLifecycleAdapter(name+"_"+clusterName, true, s, lifecycle)
	s.AddClusterScopedHandler(name, clusterName, sync)
}
