package core

import (
	"testing"

	"github.com/google/uuid"
	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/kubernetes/client"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

func Test_WatcherFetcher_Connect(t *testing.T) {
	// Must have a valid kubernetes connection

	factory := buildFactory()

	tcs := []struct {
		name        string
		cf          client.ClientFactory
		namespace   string
		expectError bool
	}{
		{
			name:        "default namespace",
			cf:          client.NewClientFactory(factory),
			namespace:   "default",
			expectError: false,
		}, {
			name:        "nil client factory",
			cf:          nil,
			namespace:   "default",
			expectError: true,
		}, {
			name:        "empty namespace",
			cf:          client.NewClientFactory(factory),
			namespace:   "",
			expectError: true,
		}, {
			name:        "garbage namespace",
			cf:          client.NewClientFactory(factory),
			namespace:   uuid.New().String(),
			expectError: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			wf := NewWatcherFetcher()
			if nil == wf {
				t.Errorf("Fatal error: `NewWatcherFetcher` returned nil")
				return
			}

			err := wf.Connect(tc.cf, tc.namespace)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error and did not receive one")
				}
			} else {
				if err != nil {
					t.Errorf("Connect returned error: %s", err.Error())
				}
			}
		})
	}
}

func Test_WatcherFetcher_Fetch(t *testing.T) {
	factory := buildFactory()

	tcs := []struct {
		name         string
		resourceType constants.ResourceType
		expectError  bool
	}{
		{
			name:         "pods",
			resourceType: constants.Pods,
			expectError:  false,
		},
		{
			name:         "empty resource type",
			resourceType: constants.ResourceType(""),
			expectError:  true,
		},
		{
			name:         "garbage resource type",
			resourceType: constants.ResourceType(uuid.New().String()),
			expectError:  true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			wf := NewWatcherFetcher()
			wf.Connect(client.NewClientFactory(factory), "default")

			wi, err := wf.Fetch(tc.resourceType)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error and did not receive one")
				}
			} else {
				if err != nil {
					t.Errorf("Connect returned error: %s", err.Error())
				}

				if wi == nil {
					t.Errorf("Got nil `watch.Interface`")
				}
			}
		})
	}
}

func buildFactory() util.Factory {
	kubeConfigFlags := genericclioptions.NewConfigFlags()
	matchVersionKubeConfigFlags := util.NewMatchVersionFlags(kubeConfigFlags)
	factory := util.NewFactory(matchVersionKubeConfigFlags)
	return factory
}
