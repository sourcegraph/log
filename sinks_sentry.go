package log

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"

	"github.com/sourcegraph/log/internal/sinkcores/sentrycore"
)

// SentrySink reports all warning-level and above log messages that contain an error field
// (via the `log.Error(err)` or `log.NamedError(name, err)` field constructors) to Sentry,
// complete with stacktrace data and any additional context logged in the corresponding
// log message (including anything accumulated on a sub-logger).
type SentrySink struct {
	// DSN configures the Sentry reporting destination.
	DSN string
	// The sample rate for event submission in the range [0.0, 1.0]. By default,
	// all events are sent. Thus, as a historical special case, the sample rate
	// 0.0 is treated as if it was 1.0. To drop all events, set the DSN to the
	// empty string.
	SampleRate float64
	// In debug mode, the debug information is printed to stdout to help you
	// understand what sentry is doing.
	Debug bool
}

type sentrySink struct {
	SentrySink

	core *sentrycore.Core
}

// NewSentrySink instantiates a Sentry sink to provide to `log.Init` with default options.
//
// See sentrycore.DefaultSentryClientOptions for the default options.
func NewSentrySink() Sink {
	return &sentrySink{SentrySink: SentrySink{}}
}

func (s *sentrySink) Name() string { return "SentrySink" }

func (s *sentrySink) build() (zapcore.Core, error) {
	client, err := sentry.NewClient(s.clientOptions())
	if err != nil {
		return nil, err
	}
	s.core = sentrycore.NewCore(sentry.NewHub(client, sentry.NewScope()))
	return s.core, nil
}

func (s *sentrySink) clientOptions() sentry.ClientOptions {
	return sentry.ClientOptions{
		Dsn:        s.DSN,
		SampleRate: s.SampleRate,
		Debug:      s.Debug,
	}
}

func (s *sentrySink) update(updated SinksConfig) error {
	if updated.Sentry == nil {
		return nil
	}

	if s.DSN == updated.Sentry.DSN && s.SampleRate == updated.Sentry.SampleRate && s.Debug == updated.Sentry.Debug {
		return nil
	}

	s.DSN = updated.Sentry.DSN
	s.SampleRate = updated.Sentry.SampleRate
	s.Debug = updated.Sentry.Debug

	client, err := sentry.NewClient(s.clientOptions())
	if err != nil {
		return err
	}

	// Do sentry setup
	s.core.SetHub(sentry.NewHub(client, sentry.NewScope()))
	return nil
}
