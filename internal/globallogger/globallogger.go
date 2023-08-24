package globallogger

import (
	"os"
	"sync"

	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/google/uuid"

	"github.com/sourcegraph/log/internal/encoders"
	"github.com/sourcegraph/log/internal/otelfields"
	"github.com/sourcegraph/log/internal/stderr"
)

var (
	EnvDevelopment = "SRC_DEVELOPMENT"
	initialized    atomic.Bool

	devMode          = os.Getenv(EnvDevelopment) == "true"
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
func Init(r otelfields.Resource, opt LoggerOption, sinks []zapcore.Core) func() error {
	globalLoggerInit.Do(func() {
		devMode = opt.Development
		globalLogger = initLogger(r, opt, sinks)
		initialized.Store(true)
	})
	return globalLogger.Sync
}

// IsInitialized indicates if the global logger is initialized.
func IsInitialized() bool { return initialized.Load() }

// forceSyncer implements the zapcore.CheckWriteHook interface and ensures that sync is called on the provided core.
// By adding it as a option with zap.WithFatalHook to the logger options, it will ensure Sync is called when a Fatal
// log is issued.
// As per the advice from https://pkg.go.dev/go.uber.org/zap#WithFatalHook, os.Exit(1) is called to halt execution after
// Sync has completed
type forceSyncer struct {
	core zapcore.Core
}

type LoggerOption struct {
	Development     bool
	StackTraceLevel zapcore.Level
}

var _ zapcore.CheckWriteHook = &forceSyncer{}

// OnWrite calls sync on the underlying core and then calls os.Exit(1).
func (f *forceSyncer) OnWrite(_ *zapcore.CheckedEntry, _ []zapcore.Field) {
	// We ignore the error here, since we're just making sure all cores have synced before exiting
	_ = f.core.Sync()
	os.Exit(1)
}

func initLogger(r otelfields.Resource, opt LoggerOption, sinks []zapcore.Core) *zap.Logger {
	internalErrsSink, err := stderr.Open()
	if err != nil {
		panic(err.Error())
	}

	options := []zap.Option{zap.ErrorOutput(internalErrsSink), zap.AddCaller()}
	if opt.Development {
		options = append(options, zap.Development())
	}

	core := zapcore.NewTee(sinks...)

	// Add a forceSyncer on the core to ensure Sync is executed on the underlying core when Fatal is called, since
	// after Fatal os.Exit is called per default configuration
	options = append(options, zap.WithFatalHook(&forceSyncer{core}))

	if opt.StackTraceLevel != zapcore.Level(0) {
		options = append(options, zap.AddStacktrace(opt.StackTraceLevel))
	}

	logger := zap.New(core, options...)

	if opt.Development {
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
