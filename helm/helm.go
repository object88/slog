package helm

import (
	"time"

	k8shelm "k8s.io/helm/pkg/helm"
)

type Helm interface {
	ListDeployments() error
}

// HelmWrapper implements the Helm interface, which abstracts away the helm
// functionality
type HelmWrapper struct {
	client  *k8shelm.Client
	timeout time.Duration
}

type Release struct {
	Name string
}

func NewHelmWrapper() *HelmWrapper {
	client := k8shelm.NewClient()

	hw := &HelmWrapper{
		client:  client,
		timeout: 30 * time.Second,
	}
	return hw
}

func (hw *HelmWrapper) ListDeployments(namespace string) ([]Release, error) {
	options := []k8shelm.ReleaseListOption{
		k8shelm.ReleaseListNamespace(namespace),
	}
	lrr, err := hw.client.ListReleases(options...)
	if err != nil {
		return nil, err
	}

	rels := lrr.Releases
	if len(rels) == 0 {
		return nil, nil
	}

	results := make([]Release, len(rels))
	for k, v := range rels {
		results[k] = Release{
			Name: v.Name,
		}
	}

	return results, nil
}
