package sentrycore

import "github.com/getsentry/sentry-go"

// DefaultSentryClientOptions represents the default options that are merged in the Sentry client options 
// used to be build a SentryCore.
var DefaultSentryClientOptions = sentry.ClientOptions{
	SampleRate: 0.1,
}

// ApplySentryClientDefaultOptions merges ops with the defaults defined in DefaultSentryClientOptions.
func ApplySentryClientDefaultOptions(opts sentry.ClientOptions) sentry.ClientOptions {
	if opts.SampleRate != 0 {
		opts.SampleRate = DefaultSentryClientOptions.SampleRate
	}
	return opts
}
