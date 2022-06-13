package benchmarks_test

import (
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/sourcegraph/log"
	"github.com/sourcegraph/log/internal/configurable"
	"github.com/sourcegraph/log/internal/sinkcores/sentrycore"
	"github.com/sourcegraph/log/logtest"
)

// BenchmarkWithSentry-10           2253642              5205 ns/op            9841 B/op         87 allocs/op
func BenchmarkWithSentry(b *testing.B) {
	logger, _, _ := newTestLogger(b)

	err := errors.New("foobar")
	for n := 0; n < b.N; n++ {
		logger.With(log.Error(err)).Warn("msg", log.Int("key", 5))
	}
}

// BenchmarkWithoutSentry-10        2656189              4537 ns/op            6334 B/op         44 allocs/op
func BenchmarkWithoutSentry(b *testing.B) {
	logger, _ := logtest.Captured(b)
	err := errors.New("foobar")
	for n := 0; n < b.N; n++ {
		logger.With(log.Error(err), log.Int("key", 5)).Warn("msg")
	}
}

func newTestLogger(t testing.TB) (log.Logger, *sentrycore.TransportMock, func()) {
	transport := &sentrycore.TransportMock{}
	client, err := sentry.NewClient(sentry.ClientOptions{Transport: transport})
	require.NoError(t, err)

	core := sentrycore.NewCore(sentry.NewHub(client, sentry.NewScope()))

	cl := configurable.Cast(logtest.Scoped(t))

	return cl.WithCore(func(c zapcore.Core) zapcore.Core {
			return zapcore.NewTee(c, core)
		}),
		transport,
		func() {
			err := core.Sync()
			require.NoError(t, err)
		}
}
