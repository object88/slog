package kubernetes

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type watcherFetcher interface {
	fetch(resource ResourceType) (watch.Interface, error)
}

type watcherFetch struct {
	client    *clientv1.CoreV1Client
	namespace string
	timeout   time.Duration
}

func (wf *watcherFetch) fetch(resource ResourceType) (watch.Interface, error) {
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
