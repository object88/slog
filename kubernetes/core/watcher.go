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
func NewWatcher() *Watcher {
	w := &Watcher{
		wf:        NewWatcherFetcher(),
		watchers:  map[constants.ResourceType]watch.Interface{},
		watchersl: &sync.Mutex{},
	}
	return w
}

func (w *Watcher) Connect(cf client.ClientFactory, namespace string) error {
	if w.wf == nil {
		return errors.New("Watcher does not have watcherFetcher; cannot connect")
	}

	return w.wf.Connect(cf, namespace)
}

// Load starts the loading procedure for the given resource type
func (w *Watcher) Load(resource constants.ResourceType) error {
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
	if wtch == nil {
		w.watchersl.Unlock()
		return errors.Errorf("Received nil watch.Interface when fetching '%s'", resource)
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
