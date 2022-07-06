package outputcore

import (
	"time"

	"github.com/sourcegraph/log/internal/encoders"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// SamplingConfig enables sampling if Initial is set. Sampling is keyed off of the log
// message only - see:
//
// - https://github.com/uber-go/zap/blob/master/FAQ.md#why-sample-application-logs
// - https://github.com/uber-go/zap/blob/master/FAQ.md#why-do-the-structured-logging-apis-take-a-message-in-addition-to-fields
type SamplingConfig struct {
	// First n entries to always log
	Initial int
	// Only log each nth entry
	Thereafter int
}

func NewCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat, sampling zap.SamplingConfig) zapcore.Core {
	core := zapcore.NewCore(
		encoders.BuildEncoder(format, false),
		output,
		level,
	)

	if sampling.Initial > 0 {
		return zapcore.NewSamplerWithOptions(core, time.Second, sampling.Initial, sampling.Thereafter)
	}
	return core
}

func NewDevelopmentCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat) zapcore.Core {
	return zapcore.NewCore(
		encoders.BuildEncoder(encoders.OutputConsole, true),
		output,
		level,
	)
}
