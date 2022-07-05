package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sourcegraph/log/internal/encoders"
	"github.com/sourcegraph/log/internal/sinkcores/outputcore"
	"github.com/sourcegraph/log/internal/stderr"
)

type outputSink struct {
	development bool

	core *outputcore.ConfigurableCore
}

func (s *outputSink) Name() string { return "OutputSink" }

func (s *outputSink) build() (zapcore.Core, error) {
	output, err := stderr.Open()
	if err != nil {
		return nil, err
	}

	level := zap.NewAtomicLevelAt(Level(os.Getenv(EnvLogLevel)).Parse())
	format := encoders.ParseOutputFormat(os.Getenv(EnvLogFormat))

	if s.development {
		return outputcore.NewDevelopmentCore(output, level, format), nil
	}

	s.core = outputcore.NewCore(output, level, format)
	return s.core, nil
}

func (s *outputSink) update(updated SinksConfig) error {
	if s.core == nil {
		return nil // not configurable
	}
	s.core.Configure(outputcore.Options{}) // TODO
	return nil
}
