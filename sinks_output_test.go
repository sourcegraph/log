package log

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/zap/zapcore"
)

func TestOutputSink_Check(t *testing.T) {
	// Unset any possible envvars that effect Check
	for _, k := range []string{
		EnvLogFormat,
		EnvLogLevel,
		EnvLogScopeLevel,
		EnvLogSamplingInitial,
		EnvLogSamplingThereafter,
	} {
		unsetenv(t, k)
	}

	logs := `
foo debug
foo error
foo.bar debug
foo.bar error
foo.bar.baz debug
foo.bar.baz error
foo.bar.baz1 debug
foo.bar.baz1 error
`

	cases := []struct {
		Name string
		Env  map[string]string
		Want string
	}{{
		Name: "error",
		Env: map[string]string{
			EnvLogLevel: "error",
		},
		Want: `
foo error
foo.bar error
foo.bar.baz error
foo.bar.baz1 error
`,
	}, {
		Name: "debug",
		Env: map[string]string{
			EnvLogLevel: "debug",
		},
		Want: logs,
	}, {
		Name: "none",
		Env: map[string]string{
			EnvLogLevel: "none",
		},
		Want: "",
	}, {
		Name: "scope",
		Env: map[string]string{
			EnvLogLevel:      "error",
			EnvLogScopeLevel: "foo.bar=debug",
		},
		// Should be everything except foo debug
		Want: `
foo error
foo.bar debug
foo.bar error
foo.bar.baz debug
foo.bar.baz error
foo.bar.baz1 debug
foo.bar.baz1 error
`,
	}, {
		Name: "scope two",
		Env: map[string]string{
			EnvLogLevel:      "error",
			EnvLogScopeLevel: "foo.bar.baz=debug,foo.bar.baz1=debug",
		},
		// Should be everything except foo and foo.bar debug
		Want: `
foo error
foo.bar error
foo.bar.baz debug
foo.bar.baz error
foo.bar.baz1 debug
foo.bar.baz1 error
`,
	}, {
		Name: "scope deep",
		Env: map[string]string{
			EnvLogLevel:      "error",
			EnvLogScopeLevel: "foo.bar.baz=debug",
		},
		Want: `
foo error
foo.bar error
foo.bar.baz debug
foo.bar.baz error
foo.bar.baz1 error
`,
	}}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			for k, v := range tc.Env {
				t.Setenv(k, v)
			}

			core, err := (&outputSink{}).build()
			if err != nil {
				t.Fatal(err)
			}

			logParts := strings.Fields(logs)
			got := ""
			for len(logParts) >= 2 {
				entry := zapcore.Entry{
					LoggerName: logParts[0],
					Level:      Level(logParts[1]).Parse(),
				}
				logParts = logParts[2:]

				if core.Check(entry, nil) != nil {
					got = fmt.Sprintf("%s%s %s\n", got, entry.LoggerName, entry.Level)
				}
			}

			want := strings.TrimSpace(tc.Want)
			got = strings.TrimSpace(got)

			if d := cmp.Diff(want, got); d != "" {
				t.Errorf("unexpected allowed logs (-want, +got):\n%s", d)
			}
		})
	}
}

func unsetenv(t *testing.T, key string) {
	v, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	os.Unsetenv(key)
	t.Cleanup(func() {
		os.Setenv(key, v)
	})
}
