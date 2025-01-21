package zLog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"testing"
)

func TestInit(t *testing.T) {

	err := Init(
		WithCommonLog("logs", "common.log", true),
		WithInfoLog("logs", "info.log", false),
		WithErrorLog("logs", "error.log", false),
		WithLogCut(1, 5, 30, true),
		WithTimeFormat("2006-01-02 15:04:05.000"))
	if err != nil {
		t.Error(err)
		return
	}
	Log.Info("hello world")
	Log.Error("hello world")
	Log.Debug("hello world")
	Log.Warn("hello world")
}

func TestDefaultLog(t *testing.T) {
	DefaultLog.Info("hello world")
	DefaultLog.Error("hello world")
	DefaultLog.Debug("hello world")
	DefaultLog.Warn("hello world")
}

func BenchmarkLog(b *testing.B) {
	err := Init(
		WithCommonLog("logs", "common.log", false), /*,
		WithInfoLog("logs", "info.log", false),
		WithErrorLog("logs", "error.log", false),
		WithLogCut(100, 5, 30, true),
		WithTimeFormat("2006-01-02 15:04:05.000")*/)
	if err != nil {
		b.Error(err)
		return
	}
	for i := 0; i < b.N; i++ {
		Log.Debug("hello world")
	}
}

func BenchmarkZap(b *testing.B) {
	_ = os.MkdirAll("logs", os.ModePerm)
	//file, _ := os.OpenFile("logs/test.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "logs/test.log",
		MaxSize:    100,
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
