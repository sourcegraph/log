package outputcore

import (
	"time"

	"github.com/sourcegraph/log/internal/encoders"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat, sampling zap.SamplingConfig, overrides []Override) zapcore.Core {
	newCore := func(level zapcore.LevelEnabler) zapcore.Core {
		return zapcore.NewCore(
			encoders.BuildEncoder(format, false),
			output,
			level,
		)
	}

	core := newOverrideCore(level, overrides, newCore)

	if sampling.Initial > 0 {
		return zapcore.NewSamplerWithOptions(core, time.Second, sampling.Initial, sampling.Thereafter)
	}
	return core
}
