package channels

import (
	"container/list"

	"k8s.io/apimachinery/pkg/watch"
)

type WatchChannel struct {
	in, out chan *watch.Event
	ll      *list.List
}

func NewWatchChannel() *WatchChannel {
	wc := &WatchChannel{
		in:  make(chan *watch.Event),
		out: make(chan *watch.Event),
		ll:  list.New(),
	}

	go wc.process()

	return wc
}

// Close satisfies io.Closer
func (wc *WatchChannel) Close() error {
	close(wc.in)
	return nil
}

func (wc *WatchChannel) In() chan<- *watch.Event {
	return wc.in
}

func (wc *WatchChannel) Out() <-chan *watch.Event {
	return wc.out
}

func (wc *WatchChannel) process() {
	in := wc.in
	var out chan *watch.Event
	var e *watch.Event
	for in != nil || out != nil {
		select {
		case v, ok := <-in:
			if !ok {
				in = nil
			} else {
				wc.ll.PushBack(v)
			}
		case out <- e:
			wc.ll.Remove(wc.ll.Front())
		}

		if wc.ll.Len() > 0 {
			out = wc.out
			e = wc.ll.Front().Value.(*watch.Event)
			// e = &e0
		} else {
			out = nil
			e = nil
		}
	}

	close(wc.out)
}
