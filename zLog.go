package zLog

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"sync"
)

type zLog struct {
	commonLogDir      string
	commonLogFileName string
	commonLogStdOut   bool

	infoLogDir      string
	infoLogFileName string
	infoLogStdOut   bool

	errorLogDir      string
	errorLogFileName string
	errorLogStdOut   bool

	maxSize    int
	maxBackups int
	maxAge     int
	compress   bool

	timeFormat string
}
type OptionFunc func(*zLog)

var once sync.Once
var logger *zLog

var Log *zap.SugaredLogger //使用前初始化

// WithCommonLog debug级别日志
func WithCommonLog(dir, fileName string, isStdOut bool) OptionFunc {
	return func(z *zLog) {
		z.commonLogDir = dir
		z.commonLogFileName = fileName
		z.commonLogStdOut = isStdOut
	}
}

// WithInfoLog info级别日志
func WithInfoLog(dir, fileName string, isStdOut bool) OptionFunc {
	return func(z *zLog) {
		z.infoLogDir = dir
		z.infoLogFileName = fileName
		z.infoLogStdOut = isStdOut
	}
}

// WithErrorLog error级别日志
func WithErrorLog(dir, fileName string, isStdOut bool) OptionFunc {
	return func(z *zLog) {
		z.errorLogDir = dir
		z.errorLogFileName = fileName
		z.errorLogStdOut = isStdOut
	}
}

// WithLogCut 日志切割
func WithLogCut(maxSize, maxBackups, maxAge int, compress bool) OptionFunc {
	return func(z *zLog) {
		z.maxSize = maxSize
		z.maxBackups = maxBackups
		z.maxAge = maxAge
		z.compress = compress
	}
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(format string) OptionFunc {
	return func(z *zLog) {
		z.timeFormat = format
	}
}

// Init 初始化
func Init(opts ...OptionFunc) error {
	var err error
	once.Do(func() {
		for _, opt := range opts {
			opt(logger)
		}
		err = initLogger()
	})
	return err
}
func initLogger() error {
	core := make([]zapcore.Core, 0)
	encoder := getEncoder()
	if logger.commonLogDir != "" && logger.commonLogFileName != "" {
		writeSyncer, err := getLogWriter(logger.commonLogDir, logger.commonLogFileName, logger.commonLogStdOut)
		if err != nil {
			return fmt.Errorf("common日志文件初始化失败:%v", err)
		}
		c1 := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
		core = append(core, c1)
	}
	if logger.errorLogDir != "" && logger.errorLogFileName != "" {
		writeSyncerErr, err := getLogWriter(logger.errorLogDir, logger.errorLogFileName, logger.errorLogStdOut)
		if err != nil {
			return fmt.Errorf("error日志文件初始化失败:%v", err)
		}
		c2 := zapcore.NewCore(encoder, writeSyncerErr, zapcore.ErrorLevel)
		core = append(core, c2)
	}
	if logger.infoLogDir != "" && logger.infoLogFileName != "" {
		writeSyncerInfo, err := getLogWriter(logger.infoLogDir, logger.infoLogFileName, logger.infoLogStdOut)
		if err != nil {
			return fmt.Errorf("info日志文件初始化失败:%v", err)
		}
		c3 := zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel)
		core = append(core, c3)
	}
	if len(core) == 0 {
		return errors.New("初始化log失败，未配置日志输出！")
	}
	// 使用NewTee将c1和c2合并到core
	multicore := zapcore.NewTee(core...)

	zapLog := zap.New(multicore, zap.AddCaller())
	Log = zapLog.Sugar()
	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if logger.timeFormat != "" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(logger.timeFormat)
	} else {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(dir, logFileName string, isStdOut bool) (zapcore.WriteSyncer, error) {
	fileName := path.Join(dir, logFileName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    logger.maxSize,
		MaxBackups: logger.maxBackups,
		MaxAge:     logger.maxAge,
		Compress:   logger.compress,
	}
	if isStdOut {
		return zapcore.AddSync(io.MultiWriter(lumberJackLogger, os.Stdout)), nil
	}
	return zapcore.AddSync(lumberJackLogger), nil
}
