package log

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hexops/autogold/v2"
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
foo info
foo error
foo.bar debug
foo.bar info
foo.bar error
foo.bar.baz debug
foo.bar.baz info
foo.bar.baz error
foo.bar.baz1 debug
foo.bar.baz1 info
foo.bar.baz1 error
`

	cases := []struct {
		Name string
		Env  map[string]string
		Want autogold.Value
	}{{
		Name: "error",
		Env: map[string]string{
			EnvLogLevel: "error",
		},
		Want: autogold.Expect(`
foo error
foo.bar error
foo.bar.baz error
foo.bar.baz1 error
`),
	}, {
		Name: "debug",
		Env: map[string]string{
			EnvLogLevel: "debug",
		},
		Want: autogold.Expect(logs),
	}, {
		Name: "none",
		Env: map[string]string{
			EnvLogLevel: "none",
		},
		Want: autogold.Expect("\n"),
	}, {
		Name: "scope",
		Env: map[string]string{
			EnvLogLevel:      "error",
			EnvLogScopeLevel: "foo.bar=debug",
		},
		// Should be everything except foo debug and info
		Want: autogold.Expect(`
foo error
foo.bar debug
foo.bar info
foo.bar error
foo.bar.baz debug
foo.bar.baz info
foo.bar.baz error
foo.bar.baz1 debug
foo.bar.baz1 info
foo.bar.baz1 error
`),
	}, {
		Name: "scope two",
		Env: map[string]string{
			EnvLogLevel:      "error",
			EnvLogScopeLevel: "foo.bar.baz=debug,foo.bar.baz1=debug",
		},
		// Should be everything except foo and foo.bar debug and info
		Want: autogold.Expect(`
foo error
foo.bar error
foo.bar.baz debug
foo.bar.baz info
foo.bar.baz error
foo.bar.baz1 debug
foo.bar.baz1 info
foo.bar.baz1 error
`),
	}, {
		Name: "scope deep",
		Env: map[string]string{
			EnvLogLevel:      "error",
			EnvLogScopeLevel: "foo.bar.baz=debug",
		},
		// Should be everything except foo and foo.bar debug and info
		Want: autogold.Expect(`
foo error
foo.bar error
foo.bar.baz debug
foo.bar.baz info
foo.bar.baz error
foo.bar.baz1 error
`),
	},
		{
			Name: "scope restricting",
			Env: map[string]string{
				EnvLogLevel:      "info",
				EnvLogScopeLevel: "foo.bar=warn",
			},
			// no debugs, and no foo.bar.* info
			Want: autogold.Expect(`
foo info
foo error
foo.bar error
foo.bar.baz error
foo.bar.baz1 error
`),
		},
	}
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

			// Add a newline to the expected output so that all expects can be on
			// newlines at the same indentation
			tc.Want.Equal(t, "\n"+got)
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
