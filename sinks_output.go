package log

import (
	"fmt"
	"os"
	"strconv"

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
		format = encoders.OutputConsole
	}

	sampling, err := parseSamplingConfig()
	if err != nil {
		return nil, err
	}
	return outputcore.NewCore(output, level, format, sampling), nil
}

// update is a no-op because outputSink cannot be changed live.
func (s *outputSink) update(updated SinksConfig) error { return nil }

func parseSamplingConfig() (config zap.SamplingConfig, err error) {
	if val, set := os.LookupEnv(EnvLogSamplingInitial); set {
		config.Initial, err = strconv.Atoi(val)
		if err != nil {
			err = fmt.Errorf("SRC_LOG_SAMPLING_INITIAL is invalid: %w", err)
			return
		}
	} else {
		config.Initial = 100
	}

	if val, set := os.LookupEnv(EnvLogSamplingInitial); set {
		config.Thereafter, err = strconv.Atoi(val)
		if err != nil {
			err = fmt.Errorf("SRC_LOG_SAMPLING_THEREAFTER is invalid: %w", err)
			return
		}
	} else {
		config.Thereafter = 100
	}

	return
}
