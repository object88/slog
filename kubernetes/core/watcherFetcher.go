package core

import (
	"time"

	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/kubernetes/client"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// WatcherFetcher gets a watch.Interface for the given resource type
type WatcherFetcher interface {
	Connect(cf client.ClientFactory, namespace string) error
	Fetch(resource constants.ResourceType) (watch.Interface, error)
}

type watcherFetch struct {
	client    clientv1.CoreV1Interface
	namespace string
	timeout   time.Duration
}

func NewWatcherFetcher() *watcherFetch {
	return &watcherFetch{
		timeout: 30 * time.Second,
	}
}

func (wf *watcherFetch) Connect(cf client.ClientFactory, namespace string) error {
	if wf == nil || cf == nil || namespace == "" {
		return errors.Errorf("Nil or invalid pointer receiver or arguments")
	}

	v1client, err := cf.V1Client()
	if err != nil {
		return err
	}

	// Validate that the namespace exists
	opts := metav1.GetOptions{}
	ns, err := v1client.Namespaces().Get(namespace, opts)
	if err != nil {
		return err
	}

	if ns == nil {
		return errors.Errorf("Nil namespace")
	}

	wf.client = v1client
	wf.namespace = namespace

	return nil
}

func (wf *watcherFetch) Fetch(resource constants.ResourceType) (watch.Interface, error) {
	if wf.client == nil {
		return nil, errors.Errorf("watcherFetch is not connected")
	}

	if err := resource.Validate(); err != nil {
		return nil, err
	}

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
