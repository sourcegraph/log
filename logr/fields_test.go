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
		"lucky_number", 3,
	})
	assert.Equal(t, fields, []log.Field{
		log.String("hello", "world"),
		log.String("goodbye", "bob"),
		log.Int("lucky_number", 3),
	})
}
