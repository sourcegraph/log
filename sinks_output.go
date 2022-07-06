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
	return outputcore.NewCore(output, level, format), nil
}

func (s *outputSink) update(updated SinksConfig) error {
	return nil
}
