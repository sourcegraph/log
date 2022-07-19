package globallogger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	assert.False(t, IsInitialized())

	// Uninitialized unsafe Get should panic
	assert.Panics(t, func() { Get(false) })

	// Uninitialized safe Get should not panic
	l, ok := Get(true)
	assert.False(t, ok)
	assert.NotNil(t, l)
}
