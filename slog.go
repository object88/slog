package slog

import (
	"fmt"

	// Ensure that OIDC is available
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/kubernetes/client"
	"github.com/object88/slog/kubernetes/core"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/watch"
	util "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

type Message struct{}
type PodStatus struct{}

type Slog struct {
	w core.Watcher

	f util.Factory

	messageOut   chan<- Message
	podStatusOut chan<- PodStatus
}

// NewSlog will return a new instance of a Slug struct
func NewSlog(factory util.Factory, messageOut chan<- Message, podStatusOut chan<- PodStatus) *Slog {
	s := &Slog{
		f:            factory,
		messageOut:   messageOut,
		podStatusOut: podStatusOut,
	}

	return s
}

// Connect will establish a RESTful connection to a Kubernetes cluster
func (s *Slog) Connect() error {
	var err error
	cf := client.NewClientFactory(s.f)
	s.w = core.NewWatcher()
	err = s.w.Connect(cf, "ecp-superquux")
	if err != nil {
		return err
	}

	return nil
}

func (s *Slog) Load(namespace string) error {
	c := s.w.Listen()
	if c == nil {
		return errors.New("Failed to get channel")
	}
	fmt.Printf("Got listener\n")
	go func(c <-chan *watch.Event) {
		for e := range c {
			switch x := e.Object.(type) {
			case *v1.Pod:
				fmt.Printf("Pod name %s: %s (%s)\n", e.Type, x.Name, x.Status.Phase)
				for _, v := range x.Status.ContainerStatuses {
					if v.State.Waiting != nil {
						fmt.Printf("  %s (%s): waiting: %s\n", v.Name, v.Image, v.State.Waiting.Reason)
					} else if v.State.Running != nil {
						fmt.Printf("  %s (%s): running since %s\n", v.Name, v.Image, v.State.Running.StartedAt.String())
					} else {
						fmt.Printf("  %s (%s): terminated: %s\n", v.Name, v.Image, v.State.Terminated.Reason)
					}
				}
			case *v1.ResourceQuota:
				fmt.Printf("Resource quota %s: %s\n", e.Type, x.Name)
			case *v1beta1.Deployment:
				fmt.Printf("Deployment name %s: %s\n", e.Type, x.Name)
			}
		}
	}(c)

	for _, rt := range constants.GetResourceTypes() {
		err := s.w.Load(rt)
		if err != nil {
			fmt.Printf("Load failed (%s): %s", rt, err.Error())
		}
	}
	fmt.Printf("Loaded\n")

	// for event := range watch.ResultChan() {
	// 	fmt.Printf("Type: %v\n", event.Type)
	// 	p, ok := event.Object.(*v1.Pod)
	// 	if !ok {
	// 		return errors.Errorf("unexpected type")
	// 	}
	// 	fmt.Printf("Statuses:\n")
	// 	for k, v := range p.Status.ContainerStatuses {
	// 		fmt.Printf("  %d: %#v\n", k, v)
	// 	}
	// 	// fmt.Printf("statuses: %#v\n", p.Status.ContainerStatuses)
	// 	fmt.Printf("phase: %s\n", p.Status.Phase)
	// }

	return nil
}
