package slog

import (
	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/kubernetes/core"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type state struct {
	w core.Watcher

	pods map[types.UID]v1.Pod
}

func newState() *state {
	return &state{}
}

func (s *state) Connect(w core.Watcher) {
	s.w = w
	go func() {
		for e := range s.w.Listen() {
			switch e.Type {
			case watch.Added:
			default:
			}
		}
	}()
}

func (s *state) Load() {
	s.w.Load(constants.Pods)
}
