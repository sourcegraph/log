package log

import "go.uber.org/zap/zapcore"

type ArrayEncoder = zapcore.ArrayEncoder
type ObjectEncoder = zapcore.ObjectEncoder

type ArrayMarshaler = zapcore.ArrayMarshaler
type ObjectMarshaler = zapcore.ObjectMarshaler

type ArrayMarshalerFunc = zapcore.ArrayMarshalerFunc
type ObjectMarshalerFunc = zapcore.ObjectMarshalerFunc
