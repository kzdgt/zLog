package zLog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"time"
)

type zLog struct {
	logDir      string               // 日志目录
	logFileName string               // 日志文件名
	isLogStdOut bool                 // 是否输出到控制台
	isJsonLog   bool                 // 是否使用json格式输出
	maxSize     int                  // 日志文件最大大小 M
	maxBackups  int                  // 日志文件最大数量
	maxAge      int                  // 日志文件最大存活时间
	compress    bool                 // 是否压缩
	level       zapcore.Level        // 日志级别
	timeFormat  string               // 时间格式
	encodeLevel zapcore.LevelEncoder // 日志级别编码器
}

type Option func(*zLog)

func (z *zLog) apply(opts ...Option) {
	for _, opt := range opts {
		opt(z)
	}
}

func New(opts ...Option) (*zap.Logger, error) {
	logger := &zLog{
		logDir:      "./logs",
		logFileName: time.Now().Format("20060102") + ".log",
		level:       zapcore.DebugLevel,
		timeFormat:  "2006-01-02 15:04:05.000",
		maxSize:     10,
		maxBackups:  3,
		maxAge:      7,
		compress:    false,
		isLogStdOut: false,
		isJsonLog:   false,
		encodeLevel: zapcore.CapitalLevelEncoder,
	}
	logger.apply(opts...)
	return logger.initLogger()
}

// WithFile 设置日志文件
func WithFile(dir, fileName string) Option {
	return func(z *zLog) {
		z.logDir = dir
		z.logFileName = fileName
	}
}

// WithLevel 设置日志级别
func WithLevel(level zapcore.Level) Option {
	return func(z *zLog) {
		z.level = level
	}
}

// WithStdOut 设置是否输出到控制台
func WithStdOut(isLogStdOut bool) Option {
	return func(z *zLog) {
		z.isLogStdOut = isLogStdOut
	}
}

// WithLogCut 日志切割
func WithLogCut(maxSize, maxBackups, maxAge int, compress bool) Option {
	return func(z *zLog) {
		z.maxSize = maxSize
		z.maxBackups = maxBackups
		z.maxAge = maxAge
		z.compress = compress
	}
}

// WithTimeFormat 设置时间格式
func WithTimeFormat(format string) Option {
	return func(z *zLog) {
		z.timeFormat = format
	}
}

// WithEncodeLevel 设置日志级别编码器
// zapcore.LowercaseColorLevelEncoder 小写颜色编码器
// zapcore.LowercaseLevelEncoder 小写编码器
// zapcore.CapitalColorLevelEncoder 大写颜色编码器
// zapcore.CapitalLevelEncoder 大写编码器
func WithEncodeLevel(levelEncoder zapcore.LevelEncoder) Option {
	return func(z *zLog) {
		z.encodeLevel = levelEncoder
	}
}

// WithJsonLog 设置是否使用json格式输出
func WithJsonLog(isJsonLog bool) Option {
	return func(z *zLog) {
		z.isJsonLog = isJsonLog
	}
}

func (z *zLog) initLogger() (*zap.Logger, error) {
	cores := make([]zapcore.Core, 0)

	logWriter, err := z.getLogWriter()
	if err != nil {
		return nil, err
	}
	cores = append(cores, zapcore.NewCore(z.getEncoder(), logWriter, z.level))

	if z.isLogStdOut {
		z.encodeLevel = zapcore.CapitalColorLevelEncoder //控制台输出增加颜色
		cores = append(cores, zapcore.NewCore(z.getEncoder(), zapcore.AddSync(os.Stdout), z.level))
	}
	core := zapcore.NewTee(cores...)
	zapLog := zap.New(core, zap.AddCaller())

	return zapLog, nil
}

func (z *zLog) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = z.encodeLevel
	if z.timeFormat != "" {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(z.timeFormat)
	} else {
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}
	if z.isJsonLog {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (z *zLog) getLogWriter() (zapcore.WriteSyncer, error) {
	fileName := path.Join(z.logDir, z.logFileName)
	if err := os.MkdirAll(z.logDir, os.ModePerm); err != nil {
		return nil, err
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    z.maxSize,
		MaxBackups: z.maxBackups,
		MaxAge:     z.maxAge,
		Compress:   z.compress,
	}

	return zapcore.AddSync(lumberJackLogger), nil
}
