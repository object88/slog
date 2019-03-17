package client

import (
	"github.com/pkg/errors"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// ClientFactory abstracts away the creation of a kubernetes core V1Client
type ClientFactory interface {
	V1Client() (clientv1.CoreV1Interface, error)
}

type clientFactory struct {
	f util.Factory
}

func NewClientFactory(f util.Factory) *clientFactory {
	return &clientFactory{
		f: f,
	}
}

func (cf *clientFactory) V1Client() (clientv1.CoreV1Interface, error) {
	if cf.f == nil {
		return nil, errors.New("Have nil factory")
	}
	restClientConfig, err := cf.f.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	v1client, err := clientv1.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return v1client, nil
}
