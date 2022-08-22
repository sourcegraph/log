package globallogger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/google/uuid"

	"github.com/sourcegraph/log/internal/encoders"
	"github.com/sourcegraph/log/internal/otelfields"
)

var (
	devMode          bool
	globalLogger     *zap.Logger
	globalLoggerInit sync.Once
)

func DevMode() bool { return devMode }

// Get retrieves the initialized global logger, or panics otherwise (unless safe is true,
// in which case a no-op logger is returned)
func Get(safe bool) *zap.Logger {
	if !IsInitialized() {
		if safe {
			return zap.NewNop()
		} else {
			panic("global logger not initialized - have you called log.Init or logtest.Init?")
		}
	}
	return globalLogger
}

// Init initializes the global logger once. Subsequent calls are no-op. Returns the
// callback to sync the root core.
func Init(r otelfields.Resource, development bool, sinks []zapcore.Core) func() error {
	globalLoggerInit.Do(func() {
		globalLogger = initLogger(r, development, sinks)
	})
	return globalLogger.Sync
}

// IsInitialized indicates if the global logger is initialized.
func IsInitialized() bool {
	return globalLogger != nil
}

// forceSyncer implements the zapcore.CheckWriteHook interface and ensures that sync is called on the provided core.
// As per the advice from https://pkg.go.dev/go.uber.org/zap#WithFatalHook, os.Exit(1) is called to halt execution after
// Sync has completed
type forceSyncer struct {
	core zapcore.Core
}

var _ zapcore.CheckWriteHook = &forceSyncer{}

// OnWrite calls sync on the underlying core and then calls os.Exit(1)
func (f *forceSyncer) OnWrite(_ *zapcore.CheckedEntry, _ []zapcore.Field) {
	f.core.Sync()
	os.Exit(1)
}

func initLogger(r otelfields.Resource, development bool, sinks []zapcore.Core) *zap.Logger {
	// Set global
	devMode = development

	internalErrsSink, err := openStderr()
	if err != nil {
		panic(err.Error())
	}

	options := []zap.Option{zap.ErrorOutput(internalErrsSink), zap.AddCaller()}
	if development {
		options = append(options, zap.Development())
	}

	core := zapcore.NewTee(sinks...)

	// Add a forceSyncer on the core to ensure Sync is executed on the underlying core when Fatal is called, since
	// after Fatal os.Exit is called per default configuration
	options = append(options, zap.WithFatalHook(&forceSyncer{core}))
	logger := zap.New(core, options...)

	if development {
		return logger
	}

	// If not in development, log OpenTelemetry Resource field and generate an InstanceID
	// to uniquely identify this resource.
	//
	// See examples: https://opentelemetry.io/docs/reference/specification/logs/data-model/#example-log-records
	if r.InstanceID == "" {
		r.InstanceID = uuid.New().String()
	}
	return logger.With(zap.Object(otelfields.ResourceFieldKey, &encoders.ResourceEncoder{Resource: r}))
}

func openStderr() (zapcore.WriteSyncer, error) {
	errSink, _, err := zap.Open("stderr")
	if err != nil {
		return nil, err
	}
	return errSink, nil
}
