// https://github.com/tjgq/broadcast/blob/master/broadcast.go
package broadcast

import (
	"sync"

	logging "github.com/ipfs/go-log"
)

var log = logging.Logger("tex-broadcast")

// Broadcaster implements a broadcast channel.
// The zero value is a usable unbuffered channel.
type Broadcaster struct {
	m         sync.Mutex
	listeners map[int]chan<- interface{} // lazy init
	nextID    int
	capacity  int
	closed    bool
}

// NewBroadcaster returns a new Broadcaster with the given capacity (0 means unbuffered).
func NewBroadcaster(n int) *Broadcaster {
	return &Broadcaster{capacity: n}
}

// Listener implements a listening endpoint for a broadcast channel.
type Listener struct {
	// Ch receives the broadcast messages.
	Ch <-chan interface{}
	b  *Broadcaster
	id int
}

// Send broadcasts a message to the channel.
func (b *Broadcaster) Send(v interface{}) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.closed {
		log.Warning("send on closed channel")
		return
	}
	for _, l := range b.listeners {
		l <- v
	}
}

// Close closes the channel, disabling the sending of further messages.
func (b *Broadcaster) Close() {
	b.m.Lock()
	defer b.m.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for _, l := range b.listeners {
		close(l)
	}
}

// Listen returns a Listener for the broadcast channel.
func (b *Broadcaster) Listen() *Listener {
	b.m.Lock()
	defer b.m.Unlock()
	if b.listeners == nil {
		b.listeners = make(map[int]chan<- interface{})
	}
	for b.listeners[b.nextID] != nil {
		b.nextID++
	}
	ch := make(chan interface{}, b.capacity)
	if b.closed {
		close(ch)
	}
	b.listeners[b.nextID] = ch
	return &Listener{ch, b, b.nextID}
}

// Close closes the Listener, disabling the receival of further messages.
func (l *Listener) Close() {
	l.b.m.Lock()
	defer l.b.m.Unlock()
	delete(l.b.listeners, l.id)
}
