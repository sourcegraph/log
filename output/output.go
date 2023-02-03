package output

// Format configures sourcegraph/log output encoding.
type Format string

const (
	// FormatJSON encodes log entries to a machine-readable, OpenTelemetry-structured
	// format.
	FormatJSON Format = "json"
	// FormatConsole encodes log entries to a human-readable format.
	FormatConsole Format = "console"
)

// ParseFormat parses the given format string as a supported output format, while
// trying to maintain some degree of back-compat with the intent of previously supported
// log formats.
func ParseFormat(format string) Format {
	switch format {
	case string(FormatJSON),
		// True 'logfmt' has significant limitations around certain field types:
		// https://github.com/jsternberg/zap-logfmt#limitations so since it implies a
		// desire for a somewhat structured format, we interpret it as OutputJSON.
		"logfmt":
		return FormatJSON
	case string(FormatConsole),
		// The previous 'condensed' format is optimized for local dev, so it serves the
		// same purpose as OutputConsole
		"condensed":
		return FormatConsole
	}

	// Fall back to JSON output
	return FormatJSON
}
