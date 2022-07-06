package log

import (
	"os"

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
	// EnvLogSamplingInitial is key of the environment variable that can be used to set
	// the number of entries with identical messages to always output per second.
	//
	// Defaults to 100 - set explicitly to 0 or -1 to disable.
	EnvLogSamplingInitial = "SRC_LOG_SAMPLING_INITIAL"
	// EnvLogSamplingThereafter is key of the environment variable that can be used to set
	// the number of entries with identical messages to discard before emitting another
	// one per second, after EnvLogSamplingInitial.
	//
	// Defaults to 100 - set explicitly to 0 or -1 to disable.
	EnvLogSamplingThereafter = "SRC_LOG_SAMPLING_THEREAFTER"
)

type Resource = otfields.Resource

// PostInitCallbacks is a set of callbacks returned by Init that enables finalization and
// updating of any configured sinks.
type PostInitCallbacks struct {
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
// Init returns a set of callbacks - see PostInitCallbacks for more details. The Sync
// callback in particular must be called before application exit.
//
// For testing, you can use 'logtest.Init' to initialize the logging library.
//
// If Init is not called, trying to create a logger with Scoped will panic.
func Init(r Resource, s ...Sink) *PostInitCallbacks {
	if globallogger.IsInitialized() {
		panic("log.Init initialized multiple times")
	}

	development := os.Getenv(EnvDevelopment) == "true"

	ss := sinks(append([]Sink{&outputSink{development: development}}, s...))
	cores, sinksBuildErr := ss.build()

	// Init the logger first, so that we can log the error if needed
	sync := globallogger.Init(r, development, cores)

	if sinksBuildErr != nil {
		// Log the error
		Scoped("log.init", "logger initialization").
			Fatal("sinks initialization failed", Error(sinksBuildErr))
	}

	return &PostInitCallbacks{
		Sync:   sync,
		Update: ss.update,
	}
}
