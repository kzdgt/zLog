# zLog

基于 [zap](https://github.com/uber-go/zap) 的高性能日志库封装，提供了更简洁的API和丰富的配置选项。

## 功能特性

- 🚀 基于 zap，性能卓越
- 📁 灵活的日志文件配置（目录、文件名）
- 🔄 自动日志切割（按大小、时间、数量）
- 🎨 支持多种日志格式（JSON、控制台）
- 🎯 可配置的日志级别
- 🕐 自定义时间格式
- 🌈 控制台彩色输出
- 📦 内置 GORM 日志适配器

## 安装

```bash
go get github.com/kzdgt/zLog
```

## 快速开始

### 基础使用

```go
package main

import (
    "github.com/kzdgt/zLog"
)

func main() {
    // 使用默认配置创建日志器
    logger, err := zLog.New()
    if err != nil {
        panic(err)
    }
    
    logger.Info("Hello, zLog!")
    logger.Error("Something went wrong")
}
```

### 自定义配置

```go
logger, err := zLog.New(
    zLog.WithFile("./logs", "app.log"),           // 设置日志文件
    zLog.WithLevel(zapcore.InfoLevel),            // 设置日志级别
    zLog.WithStdOut(true),                        // 同时输出到控制台
    zLog.WithLogCut(100, 5, 30, true),           // 日志切割配置
    zLog.WithTimeFormat("2006-01-02 15:04:05"),  // 自定义时间格式
    zLog.WithJsonLog(true),                       // 使用JSON格式
)
```

## 配置选项

### WithFile(dir, fileName string)
设置日志文件目录和文件名。

### WithLevel(level zapcore.Level)
设置日志级别，可选值：
- `zapcore.DebugLevel`
- `zapcore.InfoLevel`
- `zapcore.WarnLevel`
- `zapcore.ErrorLevel`
- `zapcore.DPanicLevel`
- `zapcore.PanicLevel`
- `zapcore.FatalLevel`

### WithStdOut(isLogStdOut bool)
设置是否同时输出到控制台，默认 `false`。

### WithLogCut(maxSize, maxBackups, maxAge int, compress bool)
配置日志切割参数：
- `maxSize`: 单个日志文件最大大小（MB）
- `maxBackups`: 保留的旧日志文件最大数量
- `maxAge`: 旧日志文件保留天数
- `compress`: 是否压缩旧日志文件

### WithTimeFormat(format string)
设置时间格式，默认 `"2006-01-02 15:04:05.000"`。

### WithEncodeLevel(encoder zapcore.LevelEncoder)
设置日志级别编码器，可选值：
- `zapcore.LowercaseLevelEncoder`: 小写
- `zapcore.LowercaseColorLevelEncoder`: 小写带颜色
- `zapcore.CapitalLevelEncoder`: 大写
- `zapcore.CapitalColorLevelEncoder`: 大写带颜色

### WithJsonLog(isJsonLog bool)
设置是否使用JSON格式输出，默认 `false`（控制台格式）。

## GORM 集成

zLog 提供了 GORM 的日志适配器：

```go
import (
    "github.com/kzdgt/zLog"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
)

// 创建zLog日志器
logger, _ := zLog.New()

// 创建GORM日志适配器
gormLogger := zLog.NewGormLogger(logger)

// 配置GORM使用zLog
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: gormLogger,
})
```

## 默认配置

```go
logDir:      "./logs"                                    // 日志目录
logFileName: time.Now().Format("20060102") + ".log"     // 日志文件名
level:       zapcore.DebugLevel                          // 日志级别
timeFormat:  "2006-01-02 15:04:05.000"                  // 时间格式
maxSize:     10                                          // 最大文件大小（MB）
maxBackups:  3                                           // 最大备份数量
maxAge:      7                                           // 最大存活天数
compress:    false                                       // 是否压缩
isLogStdOut: false                                       // 是否输出到控制台
isJsonLog:   false                                       // 是否JSON格式
encodeLevel: zapcore.CapitalLevelEncoder                 // 级别编码器
```


## 依赖

- [go.uber.org/zap](https://github.com/uber-go/zap) - 高性能日志库
- [gopkg.in/natefinch/lumberjack.v2](https://github.com/natefinch/lumberjack) - 日志文件切割

## 许可证

MIT License

## 源码地址

- GitHub: [github.com/kzdgt/zLog](https://github.com/kzdgt/zLog)