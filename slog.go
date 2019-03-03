package slog

import (
	// Ensure that OIDC is available
	"fmt"

	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	cmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

type Slog struct {
	f cmdutil.Factory

	clientset *kubernetes.Clientset

	namespace string
}

func NewSlog(factory cmdutil.Factory, service string, namespace string) *Slog {
	s := &Slog{
		f:         factory,
		namespace: namespace,
	}

	return s
}

func (s *Slog) Connect() error {
	var err error

	clientConfig := s.f.ToRawKubeConfigLoader()
	// apiConfig, err := clientConfig.RawConfig()
	// if err != nil {
	// 	return err
	// }

	restClientConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return err
	}

	s.clientset, err = kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return err
	}

	return nil
}

func (s *Slog) Load() error {
	lo := metav1.ListOptions{}

	fmt.Printf("Namespace: %s\n", s.namespace)

	pods, err := s.clientset.CoreV1().Pods(s.namespace).List(lo)
	if err != nil {
		return err
	}

	podNames := make([]string, len(pods.Items))
	for k, v := range pods.Items {
		podNames[k] = v.Name
	}

	fmt.Printf("Pod names:\n")
	for _, n := range podNames {
		fmt.Printf("%s\n", n)
	}
	return nil
}
