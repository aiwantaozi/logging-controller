package v3

import (
	"github.com/rancher/norman/lifecycle"
	"k8s.io/apimachinery/pkg/runtime"
)

type LoggingLifecycle interface {
	Create(obj *Logging) (*Logging, error)
	Remove(obj *Logging) (*Logging, error)
	Updated(obj *Logging) (*Logging, error)
}

type loggingLifecycleAdapter struct {
	lifecycle LoggingLifecycle
}

func (w *loggingLifecycleAdapter) Create(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Create(obj.(*Logging))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *loggingLifecycleAdapter) Finalize(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Remove(obj.(*Logging))
	if o == nil {
		return nil, err
	}
	return o, err
}

func (w *loggingLifecycleAdapter) Updated(obj runtime.Object) (runtime.Object, error) {
	o, err := w.lifecycle.Updated(obj.(*Logging))
	if o == nil {
		return nil, err
	}
	return o, err
}

func NewLoggingLifecycleAdapter(name string, clusterScoped bool, client LoggingInterface, l LoggingLifecycle) LoggingHandlerFunc {
	adapter := &loggingLifecycleAdapter{lifecycle: l}
	syncFn := lifecycle.NewObjectLifecycleAdapter(name, clusterScoped, adapter, client.ObjectClient())
	return func(key string, obj *Logging) error {
		if obj == nil {
			return syncFn(key, nil)
		}
		return syncFn(key, obj)
	}
}
