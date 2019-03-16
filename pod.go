package slog

import (
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Pod wraps a K8S pod object
type Pod struct {
	ID types.UID

	source *v1.Pod
}

// FromPod creates a Pod instance from a K8S pod object
func FromPod(source *v1.Pod) *Pod {
	return &Pod{
		ID:     source.GetUID(),
		source: source,
	}
}

func (p *Pod) RefreshPod(source *v1.Pod) error {
	if p.ID != source.GetUID() {
		return errors.Errorf("Internal error: mismatched pod DNS_LABEL")
	}
	p.source = source
	return nil
}

func (p *Pod) String() string {
	states := [3]int{0, 0, 0}
	for _, v := range p.source.Status.ContainerStatuses {
		if v.State.Waiting != nil {
			states[0]++
		} else if v.State.Running != nil {
			states[1]++
		} else {
			states[2]++
		}
	}

	// Ready

	// Restarts

	// Image

	return ""
}
