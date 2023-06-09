package encoders_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/hexops/autogold/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sourcegraph/log"
	"github.com/sourcegraph/log/internal/encoders"
	"github.com/sourcegraph/log/output"
)

func TestIsErrorEncoder(t *testing.T) {
	for _, f := range []log.Field{
		log.Error(errors.New("foo")),
		log.NamedError("named", errors.New("foo")),
	} {
		t.Run(f.Key, func(t *testing.T) {
			e, ok := encoders.IsErrorEncoder(f)
			assert.True(t, ok)
			assert.Error(t, e)
		})
	}

}

func TestPOC(t *testing.T) {
	buf := &bytes.Buffer{}
	l := zap.New(zapcore.NewCore(
		encoders.BuildEncoder(output.FormatJSON, encoders.EncoderOptions{
			RedactErrors: true,
		}),
		zapcore.AddSync(buf),
		zapcore.InfoLevel,
	))

	l.Error("foobar", log.Error(errors.Newf("safe: %s", "unsafe")))

	autogold.Expect(`1686282238571598000 ERROR foobar {"error": "safe: unsafe"}`).Equal(t, strings.TrimSpace(buf.String()))
}
