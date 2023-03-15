package logger

import (
	"fmt"
	"sync"

	"go.uber.org/zap/zapcore"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"go.uber.org/zap"
)

// LogWrapper wrap the logger for extension ability.
type LogWrapper struct {
	Log *zap.Logger

	event *LogEvent
	span  tracer.Span

	function string
	file     string
	line     int

	withFileTrace  bool
	withEventTrace bool
	withDDTrace    bool

	mutex sync.Mutex
}

// LogEvent holds the event information.
type LogEvent struct {
	Name   string
	Type   string
	Method string
	Path   string
}

// LogTrace holds the tracer information.
type LogTrace struct {
	SpanID  uint64
	TraceID uint64
}

// MarshalLogObject is to marshal log object.
func (l LogTrace) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddUint64("dd.trace_id", l.TraceID)
	encoder.AddUint64("dd.span_id", l.SpanID)
	return nil
}

// MarshalLogObject is to marshal log object.
func (l LogEvent) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("name", l.Name)
	encoder.AddString("type", l.Type)
	encoder.AddString("method", l.Type)
	encoder.AddString("path", l.Type)
	return nil
}

// WithEvent adds an event information about the log.
// Gives the clear information for debugging or tracing purposes.
func (l *LogWrapper) WithEvent(event *LogEvent) *LogWrapper {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.event = event

	return l
}

// WithDDTracer adds datadog tracer into log.
func (l *LogWrapper) WithDDTracer(span tracer.Span) *LogWrapper {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.span = span

	return l
}

// Info uses fmt.Sprint to construct and log a message.
func (l *LogWrapper) Info(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()

	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Info(l.buildMessage(message), l.buildFields()...)
	} else {
		l.Log.Info(l.buildMessage(message), l.buildFields()...)
	}

	return
}

// Infof uses fmt.Sprintf to construct and log a message.
func (l *LogWrapper) Infof(template string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Info(l.buildMessage(template, args...), l.buildFields()...)
	} else {
		l.Log.Info(l.buildMessage(template, args...), l.buildFields()...)

	}

	return
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *LogWrapper) Debug(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Debug(l.buildMessage(message), l.buildFields()...)
	} else {
		l.Log.Debug(l.buildMessage(message), l.buildFields()...)
	}

	return
}

// Debugf uses fmt.Sprintf to construct and log a message.
func (l *LogWrapper) Debugf(template string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Debug(l.buildMessage(template, args...), l.buildFields()...)
	} else {
		l.Log.Debug(l.buildMessage(template, args...), l.buildFields()...)
	}

	return
}

// Error uses fmt.Sprint to construct and log a message.
func (l *LogWrapper) Error(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Error(l.buildMessage(message), l.buildFields()...)
	} else {
		l.Log.Error(l.buildMessage(message), l.buildFields()...)
	}

	return
}

// Errorf uses fmt.Sprintf to construct and log a message.
func (l *LogWrapper) Errorf(template string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Error(l.buildMessage(template, args...), l.buildFields()...)
	} else {
		l.Log.Error(l.buildMessage(template, args...), l.buildFields()...)
	}

	return
}

// Fatal uses fmt.Sprint to construct and log a message.
func (l *LogWrapper) Fatal(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Fatal(l.buildMessage(message), l.buildFields()...)
	} else {
		l.Log.Fatal(l.buildMessage(message), l.buildFields()...)
	}

	return
}

// Fatalf uses fmt.Sprintf to construct and log a message.
func (l *LogWrapper) Fatalf(template string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Fatal(l.buildMessage(template, args...), l.buildFields()...)
	} else {
		l.Log.Fatal(l.buildMessage(template, args...), l.buildFields()...)
	}

	return
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *LogWrapper) Warn(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Warn(l.buildMessage(message), l.buildFields()...)
	} else {
		l.Log.Warn(l.buildMessage(message), l.buildFields()...)
	}

	return
}

// Warnf uses fmt.Sprintf to construct and log a message.
func (l *LogWrapper) Warnf(template string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Warn(l.buildMessage(template, args...), l.buildFields()...)
	} else {
		l.Log.Warn(l.buildMessage(template, args...), l.buildFields()...)
	}

	return
}

// Panic uses fmt.Sprint to construct and log a message.
func (l *LogWrapper) Panic(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Panic(l.buildMessage(message), l.buildFields()...)
	} else {
		l.Log.Panic(l.buildMessage(message), l.buildFields()...)
	}

	return
}

// Panicf uses fmt.Sprintf to construct and log a message.
func (l *LogWrapper) Panicf(template string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Panic(l.buildMessage(template, args...), l.buildFields()...)
	} else {
		l.Log.Panic(l.buildMessage(template, args...), l.buildFields()...)
	}

	return
}

// Test uses fmt.Sprint to construct and log a message.
func (l *LogWrapper) Test(template string, args ...interface{}) string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.function, l.file, l.line = getCurrentPosition()
	if l.withFileTrace {
		l.Log.
			With(
				zap.String("function", l.function),
				zap.String("file", l.file),
				zap.Int("line", l.line),
			).Info(l.buildMessage(template, args...), l.buildFields()...)
	} else {
		l.Log.Info(l.buildMessage(template, args...), l.buildFields()...)
	}

	return ""
}

func (l *LogWrapper) buildMessage(template string, fmtArgs ...interface{}) string {
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	return msg
}

func (l *LogWrapper) buildFields() []zap.Field {
	var fields []zap.Field

	if l.event != nil && l.withEventTrace {
		fields = append(fields, zap.Object("event", l.event))
	}

	if l.span != nil && l.withDDTrace {
		fields = append(fields, zap.Uint64("dd.trace_id", l.span.Context().TraceID()))
		fields = append(fields, zap.Uint64("dd.span_id", l.span.Context().SpanID()))
	}

	return fields
}
