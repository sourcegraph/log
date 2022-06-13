package log

import (
	"os"

	"go.uber.org/zap"

	"github.com/sourcegraph/log/internal/encoders"
	"github.com/sourcegraph/log/internal/globallogger"
	"github.com/sourcegraph/log/otfields"
)

const (
	envSrcDevelopment = "SRC_DEVELOPMENT"
	envSrcLogFormat   = "SRC_LOG_FORMAT"
	envSrcLogLevel    = "SRC_LOG_LEVEL"
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
// Init returns a set of callbacks - see PostInitializationCallbacks for more details.
// The Sync callback in particular must be called before application exit.
//
// For testing, you can use 'logtest.Init' to initialize the logging library.
//
// If Init is not called, Get will panic.
func Init(r Resource, s ...Sink) *PostInitializationCallbacks {
	if globallogger.IsInitialized() {
		panic("log.Init initialized multiple times")
	}

	level := zap.NewAtomicLevelAt(Level(os.Getenv(envSrcLogLevel)).Parse())
	format := encoders.ParseOutputFormat(os.Getenv(envSrcLogFormat))
	development := os.Getenv(envSrcDevelopment) == "true"

	ss := sinks(s)
	cores, err := ss.build()
	sync := globallogger.Init(r, level, format, development, cores)

	if err != nil {
		// Log the error
		Scoped("log.init", "logger initialization").Fatal("core initialization failed", Error(err))
	}

	return &PostInitializationCallbacks{
		Sync:   sync,
		Update: ss.update,
	}
}
