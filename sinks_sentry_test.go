package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSentrySink(t *testing.T) {
	t.Run("defaults are applied", func(t *testing.T) {
		s := NewSentrySink()
		ss, ok := s.(*sentrySink)
		assert.True(t, ok)
		assert.Equal(t, 0.1, ss.SentrySink.SampleRate)
	})

	t.Run("values are set", func(t *testing.T) {
		s := NewSentrySinkWith(SentrySink{DSN: "http://test:test@do.no.exist/123", Debug: true})
		ss, ok := s.(*sentrySink)
		assert.True(t, ok)
		assert.Equal(t, "http://test:test@do.no.exist/123", ss.SentrySink.DSN)
		assert.Equal(t, true, ss.SentrySink.Debug)
	})
}

func TestSentrySinkUpdate(t *testing.T) {
	t.Run("sink is updated with values", func(t *testing.T) {
		s := NewSentrySinkWith(SentrySink{DSN: "http://test:test@do.no.exist/123"})
		_, err := s.build()
		assert.Nil(t, err)

		err = s.update(SinksConfig{
			&SentrySink{
				DSN:        "",
				SampleRate: 0.3333,
			},
		})

		assert.Nil(t, err)

		ss, ok := s.(*sentrySink)
		assert.True(t, ok)

		assert.Equal(t, "", ss.DSN)
		assert.Equal(t, 0.3333, ss.SampleRate)
	})
}
