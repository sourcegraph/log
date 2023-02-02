package hook_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hexops/autogold/v2"

	"github.com/sourcegraph/log"
	"github.com/sourcegraph/log/hook"
	"github.com/sourcegraph/log/internal/encoders"
	"github.com/sourcegraph/log/logtest"
)

func TestWriter(t *testing.T) {
	logger, exportLogs := logtest.Captured(t)

	var output bytes.Buffer
	hookedLogger := hook.Writer(logger, &output, log.LevelWarn, encoders.OutputJSON)

	hookedLogger.Debug("debug")
	hookedLogger.Warn("warn")
	hookedLogger.Error("error")

	logger.Error("parent")

	// hooked logger output - only warn and above, and messages logged to parent is not
	// included
	autogold.Expect([]string{
		`{"SeverityText":"WARN","Timestamp":1675380599980541000,"InstrumentationScope":"TestWriter","Caller":"hook/writer_test.go:23","Function":"github.com/sourcegraph/log/hook_test.TestWriter","Body":"warn"}`,
		`{"SeverityText":"ERROR","Timestamp":1675380599980612000,"InstrumentationScope":"TestWriter","Caller":"hook/writer_test.go:24","Function":"github.com/sourcegraph/log/hook_test.TestWriter","Body":"error"}`,
	}).Equal(t, strings.Split(strings.TrimSpace(output.String()), "\n"))

	// parent logger output - should receive everything
	autogold.Expect([]string{"debug", "warn", "error", "parent"}).Equal(t, exportLogs().Messages())
}
