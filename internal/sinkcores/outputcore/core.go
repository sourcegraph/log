package outputcore

import (
	"github.com/sourcegraph/log/internal/encoders"
	"go.uber.org/zap/zapcore"
)

func NewCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat) zapcore.Core {
	return zapcore.NewCore(
		encoders.BuildEncoder(format, false),
		output,
		level,
	)
}

func NewDevelopmentCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat) zapcore.Core {
	return zapcore.NewCore(
		encoders.BuildEncoder(format, true),
		output,
		level,
	)
}
