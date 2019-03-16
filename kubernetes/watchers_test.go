package kubernetes

//go:generate mockgen -destination=../mocks/mock_k8s.go -package=mocks k8s.io/apimachinery/pkg/watch Interface
//go:generate mockgen -destination=../mocks/mock_util_factory.go -package=mocks k8s.io/kubernetes/pkg/kubectl/cmd/util Factory

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/object88/slog/mocks"
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

			mf := mocks.NewMockFactory(ctrl)
			mf.EXPECT().ToRESTConfig().Return(nil, nil)

			w := NewWatcher(mf, tc.namespace)
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
