package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sourcegraph/log/internal/encoders"
)

// A Field is a marshaling operation used to add a key-value pair to a logger's context.
//
// Field is an aliased import that is intentionally restricted so as to not allow overly
// liberal use of log fields, namely 'Any()'.
type Field = zapcore.Field

var (
	// String constructs a field with the given key and value.
	String = zap.String
	// Strings constructs a field that carries a slice of strings.
	Strings = zap.Strings

	// Int constructs a field with the given key and value.
	Int = zap.Int
	// Int32 constructs a field with the given key and value.
	Int32 = zap.Int32
	// Int64 constructs a field with the given key and value.
	Int64 = zap.Int64
	// Ints constructs a field that carries a slice of integers.
	Ints = zap.Ints
	// Int32s constructs a field that carries a slice of 32 bit integers.
	Int32s = zap.Int32s
	// Int64s constructs a field that carries a slice of integers.
	Int64s = zap.Int64s

	// Uint constructs a field with the given key and value.
	Uint = zap.Uint
	// Uint32 constructs a field with the given key and value.
	Uint32 = zap.Uint32
	// Uint64 constructs a field with the given key and value.
	Uint64 = zap.Uint64

	// Float32 constructs a field that carries a float32. The way the
	// floating-point value is represented is encoder-dependent, so marshaling is
	// necessarily lazy.
	Float32 = zap.Float32
	// Float32s constructs a field that carries a slice of floats.
	Float32s = zap.Float32s
	// Float64 constructs a field that carries a float64. The way the floating-point value
	// is represented is encoder-dependent, so marshaling is necessarily lazy.
	Float64 = zap.Float64
	// Float64s constructs a field that carries a slice of floats.
	Float64s = zap.Float64s

	// Bool constructs a field that carries a bool.
	Bool = zap.Bool

	// Binary constructs a field that carries an opaque binary blob.
	//
	// Binary data is serialized in an encoding-appropriate format. For example,
	// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
	// use ByteString.
	Binary = zap.Binary

	// Duration constructs a field with the given key and value. The encoder controls how
	// the duration is serialized.
	Duration = zap.Duration

	// Time constructs a Field with the given key and value. The encoder controls how the
	// time is serialized.
	Time = zap.Time

	// Namespace creates a named, isolated scope within the logger's context. All subsequent
	// fields will be added to the new namespace.
	//
	// This helps prevent key collisions when injecting loggers into sub-components or
	// third-party libraries.
	Namespace = zap.Namespace
)

// Object constructs a field that places all the given fields within the given key's
// namespace.
func Object(key string, fields ...Field) Field {
	return zap.Object(key, encoders.FieldsObjectEncoder(fields))
}

// Error is shorthand for the common idiom NamedError("error", err).
func Error(err error) Field {
	return NamedError("error", err)
}

// ptrFieldsTypes lists all acceptable pointer types for creating fields.
//
// Caveat: we can't use ~type because the type assertion switch in zap.Any will fail with a
// custom type.
type ptrFieldsTypes interface {
	*string | *int | *int32 | *int64 | *uint | *uint32 | *uint64 | *float32 | *float64 | *bool | *time.Time | *time.Duration
}

// Ptr creates a field whose value is a pointer to a type that is supported by other fields functions from this package and
// safely and explicitly represent `nil` when appropriate.
func Ptr[T ptrFieldsTypes](key string, value T) Field {
	return zap.Any(key, value)
}

// NamedError constructs a field that logs err.Error() under the provided key.
//
// For the common case in which the key is simply "error", the Error function is shorter and less repetitive.
//
// This is currently intentionally different from the zap.NamedError implementation since
// we don't want the additional verbosity at the moment.
func NamedError(key string, err error) Field {
	if err == nil {
		return String(key, "<nil>")
	}
	return zap.NamedError(key, &encoders.ErrorEncoder{Source: err})
}
