package encoders

import (
	"fmt"

	"go.uber.org/zap/zapcore"

	"github.com/cockroachdb/redact"
	"github.com/sourcegraph/log/internal/otelfields"
)

type ResourceEncoder struct {
	otelfields.Resource
}

var _ zapcore.ObjectMarshaler = &ResourceEncoder{}

func (r *ResourceEncoder) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if len(r.Name) > 0 {
		enc.AddString("service.name", r.Name)
	} else {
		enc.AddString("service.name", "unknown_service")
	}

	if len(r.Namespace) > 0 {
		enc.AddString("service.namespace", r.Namespace)
	}
	if len(r.Version) > 0 {
		enc.AddString("service.version", r.Version)
	}
	if len(r.InstanceID) > 0 {
		enc.AddString("service.instance.id", r.InstanceID)
	}
	return nil
}

type TraceContextEncoder struct{ otelfields.TraceContext }

var _ zapcore.ObjectMarshaler = &TraceContextEncoder{}

func (t *TraceContextEncoder) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if len(t.TraceID) > 0 {
		enc.AddString("TraceId", t.TraceID)
	}
	if len(t.SpanID) > 0 {
		enc.AddString("SpanId", t.SpanID)
	}
	return nil
}

type FieldsObjectEncoder []zapcore.Field

var _ zapcore.ObjectMarshaler = &FieldsObjectEncoder{}

func (fields FieldsObjectEncoder) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for _, f := range fields {
		f.AddTo(enc)
	}
	return nil
}

type ErrorEncoder struct {
	Key    string
	Source error
}

var _ zapcore.ObjectMarshaler = &FieldsObjectEncoder{}

func IsErrorEncoder(f zapcore.Field) (error, bool) {
	if f.Type != zapcore.InlineMarshalerType {
		return nil, false
	}
	if e, ok := f.Interface.(*ErrorEncoder); ok {
		return e.Source, ok
	}
	return nil, false
}

func (e *ErrorEncoder) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	fmt.Printf("%T\n", enc)

	if opts, ok := enc.(*encoderWithOptions); ok && opts.RedactErrors {
		enc.AddString(e.Key, redact.Sprintf("%v", e.Source).Redact().StripMarkers())
	} else {
		enc.AddString(e.Key, e.Source.Error())
	}
	return nil
}
