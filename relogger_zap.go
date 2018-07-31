package reloggerZap

import (
	"github.com/remicro/api/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (lgr logging.Logger, err error) {
	l, err := zap.NewProduction()
	if err != nil {
		return
	}

	lgr = &zapLogger{
		logger: l,
	}
	return
}

type zapLogger struct {
	logger *zap.Logger
	fields map[string]zap.Field
}

func (zap *zapLogger) Info() logging.Entry {
	return newEntry(zap.logger, zapcore.InfoLevel, zap.fields)
}

func (zap *zapLogger) Error() logging.Entry {
	return newEntry(zap.logger, zapcore.ErrorLevel, zap.fields)
}

func (zap *zapLogger) Debug() logging.Entry {
	return newEntry(zap.logger, zapcore.DebugLevel, zap.fields)
}

func (zap *zapLogger) Warn() logging.Entry {
	return newEntry(zap.logger, zapcore.WarnLevel, zap.fields)
}

func (zap *zapLogger) Critical() logging.Entry {
	return newEntry(zap.logger, zapcore.ErrorLevel, zap.fields)
}
