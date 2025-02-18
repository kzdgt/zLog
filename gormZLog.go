package zLog

import (
	"go.uber.org/zap"
	"time"
)

type GormLogger struct {
	zap *zap.Logger
}

func NewGormLogger(logger *zap.Logger) GormLogger {
	return GormLogger{zap: logger}
}

func (l GormLogger) Print(values ...interface{}) {
	if len(values) < 2 {
		return
	}

	switch values[0] {
	case "sql":
		l.zap.Debug("gorm.debug.sql",
			zap.String("query", values[3].(string)),
			zap.Any("values", values[4]),
			zap.Float64("duration in ms", float64(values[2].(time.Duration))/float64(time.Millisecond)),
			zap.Int64("affected-rows", values[5].(int64)),
			zap.String("source", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	default:
		l.zap.Debug("gorm.debug.other",
			zap.Any("values", values[2:]),
			zap.String("source", values[1].(string)), // if AddCallerSkip(6) is well defined, we can safely remove this field
		)
	}
}
