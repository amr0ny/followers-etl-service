package zaplogger

import (
	"github.com/amr0ny/followers-etl-service/internal/domain"
	"go.uber.org/zap"
)

type ZapLogger struct {
	sugar *zap.SugaredLogger
}

func NewZapLogger(cfg zap.Config) (domain.Logger, error) {
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{sugar: logger.Sugar()}, nil
}

func (l *ZapLogger) Info(format string, args ...any) {
	l.sugar.Infof(format, args...)
}

func (l *ZapLogger) Error(format string, args ...any) {
	l.sugar.Errorf(format, args...)
}

func (l *ZapLogger) Debug(format string, args ...any) {
	l.sugar.Debugf(format, args...)
}
