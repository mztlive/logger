package logger

import (
	"sync"

	"github.com/gookit/config/v2"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var once sync.Once

// return a zap with lumberjack logger
func createZapLogger(logOutputFile string, logDevelopment bool) *zap.Logger {
	var zapLogger *zap.Logger

	hook := lumberjack.Logger{
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,    //days
		Compress:   false, // disabled by default
		Filename:   logOutputFile,
	}

	write := zapcore.AddSync(&hook)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		write,
		zap.NewAtomicLevel(),
	)

	if logDevelopment {
		development := zap.Development()
		// zapLogger = zap.New(core, development, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
		zapLogger = zap.New(core, development, zap.AddCaller())
	} else {
		zapLogger = zap.New(core, zap.AddCaller())
		// zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	}

	return zapLogger
}

// Logger 返回一个日志组件
// 通过config读取 Logger.InfoOutputFile 和 Logger.Development
func Logger() *zap.Logger {
	once.Do(func() {
		logOutputFile := config.String("Logger.InfoOutputFile", "./logs/adsystem-crm.log")
		logDevelopment := config.Bool("Logger.Development", true)

		logger = createZapLogger(logOutputFile, logDevelopment)
	})

	return logger
}

// NewLogger returns a new logger.
// no sync.Once
func NewLogger(logOutputFile string, logDevelopment bool) *zap.Logger {
	return createZapLogger(logOutputFile, logDevelopment)
}
