package logr

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sourcegraph/log/logtest"
)

func TestGetLogger(t *testing.T) {
	logr := New(logtest.Scoped(t))

	t.Run("from the root", func(t *testing.T) {
		logger, ok := GetLogger(logr)
		assert.True(t, ok)
		assert.NotNil(t, logger)
	})

	t.Run("from a named sub-logger", func(t *testing.T) {
		logger, ok := GetLogger(logr.WithName("foobar"))
		assert.True(t, ok)
		assert.NotNil(t, logger)
	})
}
