package core

import (
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/golang/mock/gomock"
	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/mocks"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func Test_Watcher_Connect(t *testing.T) {
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

			w := NewWatcher()
			if w == nil {
				t.Errorf("Received nil watcher")
			}

			err := w.Connect(mcf, tc.namespace)
			if err != nil {
				t.Errorf("Got unexpected error: %s", err.Error())
			}
		})
	}
}

func Test_Watcher_Load(t *testing.T) {
	tcs := []struct {
		name         string
		resourceType constants.ResourceType
		expectError  bool
	}{
		{
			name:         "load deployments",
			resourceType: constants.Deployments,
		},
		{
			name:         "bad resource type",
			resourceType: constants.ResourceType("bork"),
			expectError:  true,
		},
		{
			name:         "empty resource type",
			resourceType: constants.ResourceType(""),
			expectError:  true,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mwf := mocks.NewMockWatcherFetcher(ctrl)
			mw := mocks.NewMockInterface(ctrl)

			c := make(chan watch.Event)
			defer close(c)

			if tc.expectError {
				mwf.EXPECT().Fetch(tc.resourceType).Return(nil, errors.Errorf("No")).Times(1)
			} else {
				mwf.EXPECT().Fetch(tc.resourceType).Return(mw, nil).Times(1)
				mw.EXPECT().ResultChan().Return(c).Times(1)
			}

			w := NewWatcher()
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
				return
			}

			if err != nil {
				t.Errorf("Got unexpected error from `Load`: %s", err.Error())
			}

			k8sw, ok := w.watchers[tc.resourceType]
			if !ok {
				t.Errorf("Watchers map does not include resource type")
			} else if k8sw != mw {
				t.Errorf("Watcher map has wrong element for resource type")
			}

		})
	}
}

type thing struct {
	name string
}

func newThing(name string) *thing {
	return &thing{
		name: name,
	}
}

func (t *thing) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}

func (t *thing) DeepCopyObject() runtime.Object {
	return &thing{
		name: t.name,
	}
}

func Test_Watcher_Watch(t *testing.T) {
	tcs := []struct {
		name   string
		things []runtime.Object
	}{
		{
			name:   "no things",
			things: []runtime.Object{},
		},
		{
			name:   "single thing",
			things: []runtime.Object{newThing("a")},
		},
		{
			name:   "many thing",
			things: []runtime.Object{newThing("a"), newThing("b"), newThing("c")},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			c := make(chan watch.Event)
			mw := mocks.NewMockInterface(ctrl)
			mw.EXPECT().ResultChan().Return(c).Times(1)
			mw.EXPECT().Stop().Do(func() {
				close(c)
			}).Times(1)

			mwf := mocks.NewMockWatcherFetcher(ctrl)
			mwf.EXPECT().Fetch(constants.Deployments).Return(mw, nil).Times(1)

			// Create the Watcher and fake the WatcherFetcher, so we don't have to
			// call `Connect`
			w := NewWatcher()
			w.wf = mwf

			var results []*watch.Event
			done := make(chan int)

			go func() {
				for e := range w.Listen() {
					results = append(results, e)
				}
				done <- 1
			}()

			err := w.Load(constants.Deployments)
			if err != nil {
				t.Errorf("Got unexpected error from `Load`: %s", err.Error())
			}

			// Fake some events
			for _, obj := range tc.things {
				c <- watch.Event{
					Type:   watch.Added,
					Object: obj,
				}
			}

			// Stop, and wait to be done.
			w.Stop()
			<-done

			// Check the results.
			for k, v := range results {
				if v.Object.(*thing).name != tc.things[k].(*thing).name {
					t.Errorf("Mismatched")
				}
			}
		})
	}
}
