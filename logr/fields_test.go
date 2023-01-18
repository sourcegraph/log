package logr

import (
	"testing"

	"github.com/sourcegraph/log"
	"github.com/stretchr/testify/assert"
)

func TestToLogFields(t *testing.T) {
	fields := toLogFields([]any{
		"hello", "world",
		"goodbye", "bob",
	})
	assert.Equal(t, fields, []log.Field{
		log.String("hello", "world"),
		log.String("goodbye", "bob"),
	})
}
