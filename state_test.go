package slog

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/object88/slog/internal/constants"
	"github.com/object88/slog/mocks"
)

func Test_State_SomethingSomething(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := newState()

	mw := mocks.NewMockWatcher(ctrl)
	mw.EXPECT().Load(constants.Pods)

	s.Connect(mw)

	s.Load()
}
