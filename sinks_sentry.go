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
	// DSN configures the Sentry reporting destination
	DSN string

	options sentry.ClientOptions
}

type sentrySink struct {
	SentrySink

	core *sentrycore.Core
}

// NewSentrySink instantiates a Sentry sink to provide to `log.Init` with default options.
//
// See sentrycore.DefaultSentryClientOptions for the default options.
func NewSentrySink() Sink {
	return &sentrySink{SentrySink: SentrySink{options: sentrycore.DefaultSentryClientOptions}}
}

// NewSentrySinkWithOptions instantiates a Sentry sink with advanced initial configuration
// to provide to `log.Init`. Note that configuration, notably the Sentry DSN, may be
// overwritten by subsequent calls to the `Update` callback from `log.Init`.
//
// See sentrycore.DefaultSentryClientOptions for the default options.
func NewSentrySinkWithOptions(opts sentry.ClientOptions) Sink {
	return &sentrySink{SentrySink: SentrySink{options: sentrycore.ApplySentryClientDefaultOptions(opts)}}
}

func (s *sentrySink) Name() string { return "SentrySink" }

func (s *sentrySink) build() (zapcore.Core, error) {
	opts := s.SentrySink.options
	opts.Dsn = s.DSN
	client, err := sentry.NewClient(opts)
	if err != nil {
		return nil, err
	}
	s.core = sentrycore.NewCore(sentry.NewHub(client, sentry.NewScope()))
	return s.core, nil
}

func (s *sentrySink) update(updated SinksConfig) error {
	var updatedDSN string
	if updated.Sentry != nil {
		updatedDSN = updated.Sentry.DSN
	}

	if s.DSN == updatedDSN {
		return nil
	}

	opts := s.SentrySink.options
	opts.Dsn = updatedDSN
	client, err := sentry.NewClient(opts)
	if err != nil {
		return err
	}

	// Do sentry setup
	s.core.SetHub(sentry.NewHub(client, sentry.NewScope()))
	return nil
}
