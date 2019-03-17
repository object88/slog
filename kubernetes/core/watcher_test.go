package core

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/mocks"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_Watchers_Connect(t *testing.T) {
	tcs := []struct {
		name      string
		namespace string
	}{
		{
			name:      "yes",
			namespace: "foo",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mcv1c := mocks.NewMockCoreV1Interface(ctrl)
			mcf := mocks.NewMockClientFactory(ctrl)
			mns := mocks.NewMockNamespaceInterface(ctrl)

			ns := v1.Namespace{}
			mns.EXPECT().Get(tc.namespace, gomock.Any()).Return(&ns, nil).Times(1)
			mcv1c.EXPECT().Namespaces().Return(mns).Times(1)
			mcf.EXPECT().V1Client().Return(mcv1c, nil).Times(1)

			w := NewWatcher(mcf, tc.namespace)
			if w == nil {
				t.Errorf("Received nil watcher")
			}

			err := w.Connect()
			if err != nil {
				t.Errorf("Got unexpected error: %s", err.Error())
			}
		})
	}
}

func Test_Watchers_Load(t *testing.T) {
	tcs := []struct {
		name         string
		namespace    string
		resourceType constants.ResourceType
		expectError  bool
	}{
		{
			name:         "load deployments",
			namespace:    "foo",
			resourceType: constants.Deployments,
		},
		{
			name:         "bad resource type",
			namespace:    "foo",
			resourceType: constants.ResourceType("bork"),
			expectError:  true,
		},
		{
			name:         "empty resource type",
			namespace:    "foo",
			resourceType: constants.ResourceType(""),
			expectError:  true,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mcf := mocks.NewMockClientFactory(ctrl)
			mwf := mocks.NewMockWatcherFetcher(ctrl)
			mw := mocks.NewMockInterface(ctrl)

			c := make(chan watch.Event)
			defer close(c)

			if !tc.expectError {
				mwf.EXPECT().Fetch(tc.resourceType).Return(mw, nil).Times(1)
				mw.EXPECT().ResultChan().Return(c).Times(1)
			}

			w := NewWatcher(mcf, tc.namespace)
			if w == nil {
				t.Errorf("Received nil watcher")
			}

			// Fake the WatcherFetcher, so we don't have to call `Connect`
			w.wf = mwf

			err := w.Load(tc.resourceType)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error; did not get one")
				}
			} else {
				if err != nil {
					t.Errorf("Got unexpected error from `Load`: %s", err.Error())
				}
			}

		})
	}
}