package client

import (
	"github.com/pkg/errors"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	extv1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// ClientFactory abstracts away the creation of a kubernetes core V1Client
type ClientFactory interface {
	ExtV1Client() (extv1.ExtensionsV1beta1Interface, error)
	V1Client() (clientv1.CoreV1Interface, error)
}

type clientFactory struct {
	f util.Factory
}

// NewClientFactory returns a new clientFactory instance
func NewClientFactory(f util.Factory) ClientFactory {
	return &clientFactory{
		f: f,
	}
}

func (cf *clientFactory) ExtV1Client() (extv1.ExtensionsV1beta1Interface, error) {
	restClientConfig, err := cf.restClient()
	if err != nil {
		return nil, err
	}

	extv1client, err := extv1.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return extv1client, nil
}

func (cf *clientFactory) V1Client() (clientv1.CoreV1Interface, error) {
	restClientConfig, err := cf.restClient()
	if err != nil {
		return nil, err
	}

	v1client, err := clientv1.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return v1client, nil
}

func (cf *clientFactory) restClient() (*rest.Config, error) {
	if cf.f == nil {
		return nil, errors.New("Have nil factory")
	}
	restClientConfig, err := cf.f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return restClientConfig, nil
}
