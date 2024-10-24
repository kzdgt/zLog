package zLog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"time"
)

var DefaultLog *zap.SugaredLogger

func init() {

	encoder := getDefaultEncoder()
	fileName := time.Now().Format("20060102") + ".log"
	writeSyncer := getDefaultLogWriter("./logs", fileName, true)
	c1 := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	zapLogger := zap.New(c1, zap.AddCaller())
	DefaultLog = zapLogger.Sugar()
}

func getDefaultEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000") //zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getDefaultLogWriter(dir, logFileName string, isStdOut bool) zapcore.WriteSyncer {
	fileName := path.Join(dir, logFileName)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	if isStdOut {
		return zapcore.AddSync(io.MultiWriter(lumberJackLogger, os.Stdout))
	}
	return zapcore.AddSync(lumberJackLogger)
}
