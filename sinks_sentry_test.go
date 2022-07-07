package log

import (
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"

	"github.com/sourcegraph/log/internal/sinkcores/sentrycore"
)

func TestNewSentrySink(t *testing.T) {
	t.Run("defaults are applied", func(t *testing.T) {
		s := NewSentrySink()
		ss, ok := s.(*sentrySink)
		assert.True(t, ok)
		assert.Equal(t, sentrycore.DefaultSentryClientOptions, ss.SentrySink.options)
	})
}

func TestNewSentrySinkWithOptions(t *testing.T) {
	t.Run("defaults are overridable", func(t *testing.T) {
		s := NewSentrySinkWithOptions(sentry.ClientOptions{SampleRate: 0.5})
		ss, ok := s.(*sentrySink)
		assert.True(t, ok)
		assert.Equal(t, 0.5, ss.SentrySink.options.SampleRate)
	})
	t.Run("defaults are merged", func(t *testing.T) {
		s := NewSentrySinkWithOptions(sentry.ClientOptions{ServerName: "foobar"})
		ss, ok := s.(*sentrySink)
		assert.True(t, ok)
		assert.Equal(t, "foobar", ss.SentrySink.options.ServerName)
	})
}
