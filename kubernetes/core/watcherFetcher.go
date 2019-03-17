package core

import (
	"time"

	"github.com/object88/slog/internal/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// WatcherFetcher gets a watch.Interface for the given resource type
type WatcherFetcher interface {
	Fetch(resource constants.ResourceType) (watch.Interface, error)
}

type watcherFetch struct {
	client    clientv1.CoreV1Interface
	namespace string
	timeout   time.Duration
}

func (wf *watcherFetch) Fetch(resource constants.ResourceType) (watch.Interface, error) {
	opts := metav1.ListOptions{}
	opts.Watch = true

	wtch, err := wf.client.RESTClient().
		Get().
		Namespace(wf.namespace).
		Resource(string(resource)).
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(wf.timeout).
		Watch()

	return wtch, err
}
