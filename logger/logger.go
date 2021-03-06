package logger

import (
	"fmt"
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
	args   []interface{}
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

func WithField(key string, value interface{}) *Logger {
	return logger.WithField(key, value)
}

// WithField set key and value for logger msg
func (l *Logger) WithField(key string, value interface{}) *Logger {
	ln := Logger{
		logger: l.logger,
		hooks:  l.hooks,
		args:   l.args,
	}
	ln.args = append(ln.args, key)
	ln.args = append(ln.args, value)
	return &ln
}

// Info

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func (l *Logger) Info(args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Infow(fmt.Sprint(args...), l.args...)
		return
	}
	l.logger.Info(args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Infow(fmt.Sprintf(template, args...), l.args...)
		return
	}
	l.logger.Infof(template, args...)
}

// Error

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func (l *Logger) Error(args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Errorw(fmt.Sprint(args...), l.args...)
		return
	}
	l.logger.Error(args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Errorw(fmt.Sprintf(template, args...), l.args...)
		return
	}
	l.logger.Errorf(template, args...)
}

// Debug

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func (l *Logger) Debug(args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Debugw(fmt.Sprint(args...), l.args...)
		return
	}
	l.logger.Debug(args...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Debugw(fmt.Sprintf(template, args...), l.args...)
		return
	}
	l.logger.Debugf(template, args...)
}

// Warn

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Warnw(fmt.Sprint(args...), l.args...)
		return
	}
	l.logger.Warn(args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Warnw(fmt.Sprintf(template, args...), l.args...)
		return
	}
	l.logger.Warnf(template, args...)
}

// Fatal

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Fatalw(fmt.Sprint(args...), l.args...)
		return
	}
	l.logger.Fatal(args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Fatalw(fmt.Sprintf(template, args...), l.args...)
		return
	}
	l.logger.Fatalf(template, args...)
}

// Panic

func Panic(args ...interface{}) {
	logger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Panicf(template, args...)
}

func (l *Logger) Panic(args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Panicw(fmt.Sprint(args...), l.args...)
		return
	}
	l.logger.Panic(args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	if len(l.args) != 0 {
		l.logger.Panicw(fmt.Sprintf(template, args...), l.args...)
		return
	}
	l.logger.Panicf(template, args...)
}
