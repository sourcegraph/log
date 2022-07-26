package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sourcegraph/log"
	"github.com/sourcegraph/log/internal/globallogger"
	"github.com/sourcegraph/log/internal/otelfields"
	"github.com/sourcegraph/log/logtest"
)

func TestLogger(t *testing.T) {
	logger, exportLogs := logtest.Captured(t)
	assert.NotNil(t, logger)

	// HACK: If in devmode, the attributes namespace does not get added, but we want to
	// test that behaviour here so we add it back.
	if globallogger.DevMode() {
		logger = logger.With(otelfields.AttributesNamespace)
	}

	logger.Debug("a debug message") // 0

	logger = logger.With(log.String("some", "field"))

	logger.Info("hello world", log.String("hello", "world")) // 1

	logger = logger.WithTrace(log.TraceContext{TraceID: "1234abcde"})
	logger.Info("goodbye", log.String("world", "hello")) // 2
	logger.Warn("another message")                       // 3

	logger.Error("object of fields", // 4
		log.Object("object",
			log.String("field1", "value"),
			log.String("field2", "value"),
		))

	logs := exportLogs()
	assert.Len(t, logs, 5)
	for _, l := range logs {
		assert.Equal(t, "TestLogger", l.Scope) // scope is always applied
	}

	// Nested fields should be in attributes
	assert.Equal(t, map[string]interface{}{
		"some":  "field",
		"hello": "world",
	}, logs[1].Fields["Attributes"])

	// TraceId should be in root, everything else in attributes
	assert.Equal(t, "1234abcde", logs[2].Fields["TraceId"])
	assert.Equal(t, map[string]interface{}{
		"some":  "field",
		"world": "hello",
	}, logs[2].Fields["Attributes"])

	// Nested fields should be in attributes
	assert.Equal(t, map[string]interface{}{
		"some": "field",
		"object": map[string]interface{}{
			"field1": "value",
			"field2": "value",
		},
	}, logs[4].Fields["Attributes"])
}
