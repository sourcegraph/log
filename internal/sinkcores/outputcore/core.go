package outputcore

import (
	"time"

	"github.com/sourcegraph/log/internal/encoders"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
