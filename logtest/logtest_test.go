package logtest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sourcegraph/log"
)

func TestExport(t *testing.T) {
	logger, exportLogs := Captured(t)
	assert.NotNil(t, logger)

	logger.Info("hello world", log.String("key", "value"))
	logger.Error("goodbye world", log.String("key", "value"))

	logs := exportLogs()
	assert.Len(t, logs, 2)
	assert.Equal(t, "TestExport", logs[0].Scope)    // test name is the scope
	assert.Equal(t, "hello world", logs[0].Message) // retains the message

	// In dev mode, attributes are not added, but custom fields are retained
	assert.Equal(t, map[string]interface{}{"key": "value"}, logs[0].Fields)

	// We can filter for entries
	assert.Len(t, logs.Filter(func(l CapturedLog) bool {
		return l.Level == log.LevelError
	}), 1)

	// We can assert the existence of an entry
	assert.False(t, logs.Contains(func(l CapturedLog) bool {
		return l.Level == log.LevelWarn
	}))
}
