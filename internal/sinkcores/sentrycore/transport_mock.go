package sentrycore

import (
	"context"
	"sync"
	"time"

	"github.com/getsentry/sentry-go"
)

type TransportMock struct {
	mu        sync.Mutex
	events    []*sentry.Event
	lastEvent *sentry.Event
}

func (t *TransportMock) Configure(options sentry.ClientOptions) {}

func (t *TransportMock) SendEvent(event *sentry.Event) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, event)
	t.lastEvent = event
}

func (t *TransportMock) Flush(timeout time.Duration) bool {
	return true
}

func (t *TransportMock) FlushWithContext(ctx context.Context) bool {
	return true
}

func (t *TransportMock) Close() {}

func (t *TransportMock) Events() []*sentry.Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.events
}
