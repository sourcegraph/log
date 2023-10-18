package main

import (
	"os"
	"time"

	"github.com/sourcegraph/log"
)

func main() {
	liblog := log.Init(log.Resource{
		Name: "logexample",
	})
	defer liblog.Sync()

	l := log.Scoped("foo")

	// print diagnostics
	config := []log.Field{}
	for _, k := range []string{
		log.EnvDevelopment,
		log.EnvLogFormat,
		log.EnvLogLevel,
		log.EnvLogScopeLevel,
		log.EnvLogSamplingInitial,
		log.EnvLogSamplingThereafter,
	} {
		config = append(config, log.String(k, os.Getenv(k)))
	}
	l.Info("configuration", config...)

	// sample message
	l.Warn("hello world!", log.Time("now", time.Now()))
}
