package outputcore

import (
	"strings"
	"time"

	"github.com/sourcegraph/log/internal/encoders"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Override allows adjusting the log level for specific scopes.
type Override struct {
	// Scope matches any zapcore.Entry with the name Scope or the prefix Scope + ".".
	Scope string
	// Level is the level to log at for zapcore.Entry's that match Scope.
	Level zapcore.Level
}

func NewCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat, sampling zap.SamplingConfig, overrides []Override) zapcore.Core {
	minOverrideLevel := level
	for _, o := range overrides {
		if !minOverrideLevel.Enabled(o.Level) {
			minOverrideLevel = o.Level
		}
	}

	core := &overrideCore{
		enc: encoders.BuildEncoder(format, false),
		out: output,

		level:            level,
		minOverrideLevel: minOverrideLevel,

		overrides: overrides,
	}

	if sampling.Initial > 0 {
		return zapcore.NewSamplerWithOptions(core, time.Second, sampling.Initial, sampling.Thereafter)
	}
	return core
}

// override is a zapcore.Core and zapcore.LevelEnabler which additionally
// will logs entries were allow returns true no matter the level.
//
// This does not wrap another core since that would not allow us to override
// that core's level. It is mostly copy pasta of zapcore.ioCore
type overrideCore struct {
	enc zapcore.Encoder
	out zapcore.WriteSyncer

	// level is the passed in level
	level zapcore.LevelEnabler

	// minOverrideLevel is the most verbose level we could log at. We have to
	// return this level as part of our interface since if we have any child
	// cores they can check directly against us via Enabled.
	minOverrideLevel zapcore.LevelEnabler

	overrides []Override
}

func (c *overrideCore) Enabled(l zapcore.Level) bool {
	return c.minOverrideLevel.Enabled(l)
}

func (c *overrideCore) Level() zapcore.Level {
	return zapcore.LevelOf(c.minOverrideLevel)
}

func (c *overrideCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	for i := range fields {
		fields[i].AddTo(c.enc)
	}
	return clone
}

func (c *overrideCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.level.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	if !c.minOverrideLevel.Enabled(ent.Level) {
		return ce
	}

	for _, o := range c.overrides {
		if !strings.HasPrefix(ent.LoggerName, o.Scope) {
			continue
		}
		// Check that if o.Scope != ent.LoggerName then it is a child logger
		if len(ent.LoggerName) > len(o.Scope) && ent.LoggerName[len(o.Scope)] != '.' {
			continue
		}

		if o.Level.Enabled(ent.Level) {
			return ce.AddCore(ent, c)
		}
	}

	return ce
}

func (c *overrideCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	_, err = c.out.Write(buf.Bytes())
	buf.Free()
	if err != nil {
		return err
	}
	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		c.Sync()
	}
	return nil
}

func (c *overrideCore) Sync() error {
	return c.out.Sync()
}

func (c *overrideCore) clone() *overrideCore {
	return &overrideCore{
		enc: c.enc.Clone(),
		out: c.out,

		level:            c.level,
		minOverrideLevel: c.minOverrideLevel,

		overrides: c.overrides,
	}
}
