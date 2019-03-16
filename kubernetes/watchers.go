package kubernetes

import (
	"sync"
	"time"

	// Ensure that OIDC is available
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// ResourceType describes the K8S types available
type ResourceType string

const (
	Deployments ResourceType = "deployments"
	Pods        ResourceType = "pods"
)

// Watcher observes K8S resources
type Watcher struct {
	namespace string
	timeout   time.Duration

	watchers  map[ResourceType]watch.Interface
	watchersl sync.Locker

	f  util.Factory
	wf watcherFetcher

	// clientset        *kubernetes.Clientset
	// restClientConfig *rest.Config
}

// NewWatcher returns a new instance of a watcher struct.  The inputs are not
// validated
func NewWatcher(factory util.Factory, namespace string) *Watcher {
	w := &Watcher{
		f:         factory,
		namespace: namespace,
		watchers:  map[ResourceType]watch.Interface{},
		watchersl: &sync.Mutex{},
	}
	return w
}

func (w *Watcher) Connect() error {
	if w.f == nil {
		return errors.Errorf("No factory")
	}

	// clientConfig := w.f.ToRawKubeConfigLoader()
	// apiConfig, err := clientConfig.RawConfig()
	// if err != nil {
	// 	return err
	// }

	restClientConfig, err := w.f.ToRESTConfig() // clientConfig.ClientConfig()
	if err != nil {
		return err
	}

	v1client, err := clientv1.NewForConfig(restClientConfig)
	if err != nil {
		return err
	}

	w.wf = &watcherFetch{
		client:    v1client,
		namespace: w.namespace,
		timeout:   30 * time.Second,
	}

	// w.clientset, err = kubernetes.NewForConfig(restClientConfig)
	// if err != nil {
	// 	return err
	// }

	// Validate that the namespace exists
	opts := metav1.GetOptions{}
	_, err = v1client.Namespaces().Get(w.namespace, opts)
	if err != nil {
		return err
	}

	return nil
}

// Load starts the loading procedure for the given resource type
func (w *Watcher) Load(resource ResourceType) error {
	w.watchersl.Lock()

	if _, ok := w.watchers[resource]; ok {
		w.watchersl.Unlock()
		return nil
	}

	wtch, err := w.wf.fetch(resource)
	if err != nil {
		w.watchersl.Unlock()
		return err
	}

	c := wtch.ResultChan()

	go func(watchChan <-chan watch.Event) {
		for v := range watchChan {
			_, ok := v.Object.(*v1.ResourceQuota)
			if !ok {
				continue
			}

			switch v.Type {
			case watch.Added:

			case watch.Deleted:

			case watch.Modified:

			case watch.Error:

			}
		}
	}(c)

	w.watchers[resource] = wtch

	w.watchersl.Unlock()

	return nil
}

func (w *Watcher) Stop() {
	w.watchersl.Lock()

	// Stop all the watchers
	for _, v := range w.watchers {
		v.Stop()
	}

	// Be sure to clear the map
	w.watchers = map[ResourceType]watch.Interface{}

	w.watchersl.Unlock()
}
