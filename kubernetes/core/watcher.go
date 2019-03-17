package core

import (
	"sync"
	"time"

	// Ensure that OIDC is available
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/kubernetes/client"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// Watcher observes K8S resources
type Watcher struct {
	namespace string
	timeout   time.Duration

	watchers  map[constants.ResourceType]watch.Interface
	watchersl sync.Locker

	cf client.ClientFactory
	wf WatcherFetcher
}

// NewWatcher returns a new instance of a watcher struct.  The inputs are not
// validated
func NewWatcher(cf client.ClientFactory, namespace string) *Watcher {
	w := &Watcher{
		cf:        cf,
		namespace: namespace,
		watchers:  map[constants.ResourceType]watch.Interface{},
		watchersl: &sync.Mutex{},
	}
	return w
}

func (w *Watcher) Connect() error {
	if w.cf == nil {
		return errors.New("Watcher does not have clientFactory; cannot connect")
	}

	v1client, err := w.cf.V1Client()
	if err != nil {
		return err
	}

	w.wf = &watcherFetch{
		client:    v1client,
		namespace: w.namespace,
		timeout:   30 * time.Second,
	}

	// Validate that the namespace exists
	opts := metav1.GetOptions{}
	ns, err := v1client.Namespaces().Get(w.namespace, opts)
	if err != nil {
		return err
	}

	if ns == nil {
		return errors.Errorf("Could not find namespace with name '%s'", w.namespace)
	}

	return nil
}

// Load starts the loading procedure for the given resource type
func (w *Watcher) Load(resource constants.ResourceType) error {
	if err := resource.Validate(); err != nil {
		return err
	}

	w.watchersl.Lock()

	if _, ok := w.watchers[resource]; ok {
		w.watchersl.Unlock()
		return nil
	}

	wtch, err := w.wf.Fetch(resource)
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
	w.watchers = map[constants.ResourceType]watch.Interface{}

	w.watchersl.Unlock()
}
