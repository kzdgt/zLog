package zLog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	logger, err := New(
		WithFile("logs", "common.log"),
		WithLogCut(1, 5, 30, true),
		WithTimeFormat("2006-01-02 15:04:05.000"))
	if err != nil {
		t.Error(err)
		return
	}
	logger.Info("hello world")
	logger.Error("hello world")
	logger.Debug("hello world")
	logger.Warn("hello world")
}

func TestDefaultLog(t *testing.T) {
	logger, err := New()
	if err != nil {
		t.Error(err)
		return
	}
	logger.Info("hello world")
	logger.Error("hello world")
	logger.Debug("hello world")
	logger.Warn("hello world")
}

func BenchmarkLog(b *testing.B) {
	logger, err := New(
		WithFile("logs", "common.log"),
		WithLogCut(1, 5, 30, false),
		WithTimeFormat("2006-01-02 15:04:05.000"))
	if err != nil {
		b.Error(err)
		return
	}
	slog := logger.Sugar()
	for i := 0; i < b.N; i++ {
		slog.Debug("hello world")
	}
}

func BenchmarkZap(b *testing.B) {
	_ = os.MkdirAll("logs", os.ModePerm)
	//file, _ := os.OpenFile("logs/test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs/test.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	writeSyncer := zapcore.AddSync(lumberJackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), writeSyncer, zapcore.DebugLevel)

	logger := zap.New(core, zap.AddCaller())
	sugarLogger := logger.Sugar()
	for i := 0; i < b.N; i++ {
		sugarLogger.Debug("hello world")
	}
}
