package slog

import (
	// Ensure that OIDC is available
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"github.com/object88/slog/kubernetes/client"
	"github.com/object88/slog/kubernetes/core"
	util "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

type Message struct{}
type PodStatus struct{}

type Slog struct {
	w *core.Watcher

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
	// lo := metav1.ListOptions{}

	// pods := s.clientset.CoreV1().Pods(namespace)

	// podList, err := pods.List(lo)
	// if err != nil {
	// 	return err
	// }

	// podNames := make([]string, len(podList.Items))
	// for k, v := range podList.Items {
	// 	podNames[k] = v.Name
	// }

	// fmt.Printf("Pod names:\n")
	// for _, n := range podNames {
	// 	fmt.Printf("  %s\n", n)
	// }

	// emc, err := external_metrics.NewForConfig(s.restClientConfig)
	// if err != nil {
	// 	return err
	// }

	// // Want to verify that this resource exists
	// mi := emc.NamespacedMetrics(namespace)
	// _, err = mi.List("*", labels.Everything())
	// if err != nil {
	// 	fmt.Printf("Nope: %s\n", err.Error())
	// 	return err
	// }

	// watch, err := pods.Watch(metav1.ListOptions{
	// 	LabelSelector: "",
	// })
	// if err != nil {
	// 	return err
	// }
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
