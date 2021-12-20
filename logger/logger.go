package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var logger *Logger

type Logger struct {
	logger *zap.SugaredLogger
	hooks  *lumberjack.Logger
}

type Options struct {
	Level      string
	Json       bool   // 是否 json 格式
	Filename   string //日志文件位置
	MaxSize    int    // 单文件最大容量,单位是MB
	MaxBackups int    // 最大保留过期文件个数
	MaxAge     int    // 保留过期文件的最大时间间隔,单位是天
	Compress   bool   // 是否需要压缩滚动日志, 使用的 gzip 压缩
}

func InitLogger() {
	logger = GetLogger()
}

func InitLoggerByOptions(options *Options) {
	logger = GetLoggerByOptions(options)
}

func Close() error {
	return logger.Close()
}

func GetLogger() *Logger {
	return GetLoggerByOptions(&Options{Level: "debug"})
}

func GetLoggerByOptions(options *Options) *Logger {
	logger := &Logger{}
	logger.hooks = options.getLumberjack()
	logger.initLogger(options)
	return logger
}

func (l *Logger) Close() error {
	if l.hooks != nil {
		if err := l.hooks.Close(); err != nil {
			return err
		}
	}
	if err := l.logger.Sync(); err != nil {
		return err
	}
	return nil
}

func (l *Logger) initLogger(options *Options) {
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	switch strings.ToLower(options.Level) {
	case "fatal":
		atomicLevel.SetLevel(zap.FatalLevel)
	case "panic":
		atomicLevel.SetLevel(zap.PanicLevel)
	case "error":
		atomicLevel.SetLevel(zap.ErrorLevel)
	case "warn":
		atomicLevel.SetLevel(zap.WarnLevel)
	case "info":
		atomicLevel.SetLevel(zap.InfoLevel)
	case "debug":
		atomicLevel.SetLevel(zap.DebugLevel)
	default:
		atomicLevel.SetLevel(zap.DebugLevel)
	}
	// 公用编码器
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
	var encoder zapcore.Encoder
	if options.Json {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	// 打印
	var writeSyncer zapcore.WriteSyncer
	if l.hooks != nil {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(l.hooks))
	} else {
		writeSyncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		encoder,     // 编码器配置
		writeSyncer, // 打印
		atomicLevel, // 日志级别
	)
	l.logger = zap.New(core).Sugar()
}

func (o *Options) getLumberjack() *lumberjack.Logger {
	filename := o.Filename
	maxSize := o.MaxSize
	maxBackup := o.MaxBackups
	maxAge := o.MaxAge
	compress := o.Compress

	if filename == "" {
		return nil
	}
	if maxSize == 0 {
		maxSize = 1
	}
	if maxBackup == 0 {
		maxBackup = 10
	}
	if maxAge == 0 {
		maxAge = 10
	}

	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}

// Info

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

// Error

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Debug

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Warn

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

// Fatal

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

// Panic

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Panicw(msg, keysAndValues...)
}

func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}
