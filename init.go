package log

import (
	"os"

	"go.uber.org/zap"

	"github.com/sourcegraph/log/internal/encoders"
	"github.com/sourcegraph/log/internal/globallogger"
	"github.com/sourcegraph/log/otfields"
)

const (
	// EnvDevelopment is key of the environment variable that is used to set whether
	// to use development logger configuration on Init.
	EnvDevelopment = "SRC_DEVELOPMENT"
	// EnvLogLevel is key of the environment variable that is used to set the log format
	// on Init.
	EnvLogFormat = "SRC_LOG_FORMAT"
	// EnvLogLevel is key of the environment variable that can be used to set the log
	// level on Init.
	EnvLogLevel = "SRC_LOG_LEVEL"
)

type Resource = otfields.Resource

// PostInitializationCallbacks wraps the callbacks that enables to sync and update the
// sinks used by the logger on configuration changes.
type PostInitializationCallbacks struct {
	// Sync must be called before application exit, such as via defer.
	Sync func() error

	// Update should be called to change sink configuration, e.g. via
	// conf.Watch. Note that sinks not created upon initialization will
	// not be created post-initialization. Is a no-op if no sinks are enabled.
	Update func(SinksConfigGetter) func()
}

// Init initializes the log package's global logger as a logger of the given resource.
// It must be called on service startup, i.e. 'main()', NOT on an 'init()' function.
// Subsequent calls will panic, so do not call this within a non-service context.
//
// Init returns a callback, sync, that should be called before application exit.
//
// For testing, you can use 'logtest.Init' to initialize the logging library.
//
// If Init is not called, Get will panic.
func Init(r Resource, s ...Sink) *PostInitializationCallbacks {
	if globallogger.IsInitialized() {
		panic("log.Init initialized multiple times")
	}

	level := zap.NewAtomicLevelAt(Level(os.Getenv(EnvLogLevel)).Parse())
	format := encoders.ParseOutputFormat(os.Getenv(EnvLogFormat))
	development := os.Getenv(EnvDevelopment) == "true"

	sinks := Sinks(s)
	update := sinks.Update
	cores, err := sinks.Build()
	sync := globallogger.Init(r, level, format, development, cores)

	if err != nil {
		// Log the error
		Scoped("log.init", "logger initialization").Fatal("core initialization failed", Error(err))
	}

	return &PostInitializationCallbacks{
		Sync:   sync,
		Update: update,
	}
}
