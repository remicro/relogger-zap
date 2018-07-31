package reloggerZap

import (
	"fmt"
	"github.com/remicro/api/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type zapFields map[string]zap.Field

func newEntry(logger *zap.Logger, level zapcore.Level, fields zapFields) *entry {
	if fields == nil {
		fields = make(zapFields)
	}
	return &entry{
		logger: logger,
		level:  level,
		fields: fields,
	}
}

type entry struct {
	level  zapcore.Level
	fields zapFields
	logger *zap.Logger
}

func (e *entry) String(key, value string) logging.Entry {
	e.fields[key] = zap.String(key, value)
	return e
}

func (e *entry) Int(key string, value int) logging.Entry {
	e.fields[key] = zap.Int64(key, int64(value))
	return e
}

func (e *entry) Err(err error) logging.Entry {
	e.fields["error"] = zap.Error(err)
	return e
}

func (e *entry) Bool(key string, value bool) logging.Entry {
	e.fields[key] = zap.Bool(key, value)
	return e
}

func (e *entry) Time(key string, value time.Time) logging.Entry {
	e.fields[key] = zap.Time(key, value)
	return e
}

func (e *entry) Duration(key string, duration time.Duration) logging.Entry {
	e.fields[key] = zap.Duration(key, duration)
	return e
}

func (e *entry) Float64(key string, value float64) logging.Entry {
	e.fields[key] = zap.Float64(key, value)
	return e
}

func (e *entry) Uint64(key string, value uint64) logging.Entry {
	e.fields[key] = zap.Uint64(key, value)
	return e
}

func (e *entry) Logf(message string, args ...interface{}) {
	e.Log(fmt.Sprintf(message, args...))
}

func (e *entry) Log(msg string) {
	defer e.logger.Sync()
	var fields []zap.Field

	if len(e.fields) > 0 {
		fields = make([]zap.Field, 0, len(e.fields))
		for _, v := range e.fields {
			fields = append(fields, v)
		}
	}

	switch e.level {
	case zapcore.DebugLevel:
		e.logger.Debug(msg, fields...)
	case zapcore.ErrorLevel:
		e.logger.Error(msg, fields...)
	case zapcore.WarnLevel:
		e.logger.Warn(msg, fields...)
	default:
		e.logger.Info(msg, fields...)
	}
}
