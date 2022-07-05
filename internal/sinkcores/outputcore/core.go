package outputcore

import (
	"sync"
	"time"

	"github.com/sourcegraph/log/internal/encoders"
	"go.uber.org/zap/zapcore"
)

func NewCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat) *ConfigurableCore {
	return newConfigurable(zapcore.NewCore(
		encoders.BuildEncoder(format, false),
		output,
		level,
	))
}

func NewDevelopmentCore(output zapcore.WriteSyncer, level zapcore.LevelEnabler, format encoders.OutputFormat) zapcore.Core {
	return zapcore.NewCore(
		encoders.BuildEncoder(format, true),
		output,
		level,
	)
}

func newConfigurable(core zapcore.Core) *ConfigurableCore {
	c := &ConfigurableCore{
		destination: core,
	}
	c.Configure(Options{})
	return c
}

type ConfigurableCore struct {
	mux sync.RWMutex

	// wrapped wraps destination, and can be rebuilt
	wrapped     zapcore.Core
	destination zapcore.Core
}

type Options struct {
	SamplingFirst      int
	SamplingThereafter int
}

func (s *ConfigurableCore) Configure(opts Options) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if opts.SamplingFirst > 0 {
		s.wrapped = zapcore.NewSamplerWithOptions(s.destination, time.Second, opts.SamplingFirst, opts.SamplingThereafter)
	} else {
		s.wrapped = s.destination // otherwise use the underlying core directly
	}
}

// Implement the zapcore.Core interface
var _ zapcore.Core = &ConfigurableCore{}

func (s *ConfigurableCore) Enabled(l zapcore.Level) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.wrapped.Enabled(l)
}

func (s *ConfigurableCore) With(f []zapcore.Field) zapcore.Core {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.wrapped.With(f)
}

func (s *ConfigurableCore) Check(e zapcore.Entry, c *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.wrapped.Check(e, c)
}

func (s *ConfigurableCore) Write(e zapcore.Entry, f []zapcore.Field) error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.wrapped.Write(e, f)
}

// Sync flushes buffered logs (if any).
func (s *ConfigurableCore) Sync() error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.destination.Sync()
}
