package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	gLogger "gorm.io/gorm/logger"
	gUtils "gorm.io/gorm/utils"
)

var zlog *zap.SugaredLogger
var once sync.Once

const (
	ColorRed     = "\033[31m"
	ColorReset   = "\033[0m"
	LogFormatStr = "LOG_FORMAT"
	JsonMode     = "JSON"
)

var logFormat string

func init() {
	once.Do(func() {
		logFormat = os.Getenv(LogFormatStr)
		zlog = NewLogger()
		zlog = zlog.WithOptions(zap.AddCallerSkip(1))
	})
}

func NewLogger() *zap.SugaredLogger {
	atom := zap.NewAtomicLevel()
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "warning":
		atom.SetLevel(zap.WarnLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	case "fatal":
		atom.SetLevel(zap.FatalLevel)
	default:
		atom.SetLevel(zap.InfoLevel)
	}

	var cfg zap.Config
	if logFormat == JsonMode {
		cfg = zap.NewProductionConfig() // JSON format
	} else {
		cfg = zap.NewDevelopmentConfig()                                 // Plain text format
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // Color output for better readability in terminal
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	}
	cfg.Level = atom
	logger := zap.Must(cfg.Build())
	sugaredLogger := logger.Sugar()
	return sugaredLogger
}

func Debugf(format string, args ...interface{}) {
	zlog.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	zlog.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	zlog.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	zlog.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	zlog.Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	zlog.Debug(args...)
}

func Info(args ...interface{}) {
	zlog.Info(args...)
}

func Warn(args ...interface{}) {
	zlog.Warn(args...)
}

func Error(args ...interface{}) {
	zlog.Error(args...)
}

func Fatal(args ...interface{}) {
	zlog.Fatal(args...)
}

func Debugln(args ...interface{}) {
	zlog.Debugln(args...)
}

func Infoln(args ...interface{}) {
	zlog.Infoln(args...)
}

func Warnln(args ...interface{}) {
	zlog.Warnln(args...)
}

func Errorln(args ...interface{}) {
	zlog.Errorln(args...)
}

func Fatalln(args ...interface{}) {
	zlog.Fatalln(args...)
}

var _ io.Writer = &Log{}
var _ gLogger.Interface = &Log{}

type Log struct {
	// LogLevel determines the minimum log level that will be recorded.
	// It uses the gLogger.LogLevel enum to define levels such as Info, Warn, Error, etc.
	LogLevel gLogger.LogLevel
}

func DefualtLog() *Log {
	return &Log{
		LogLevel: gLogger.Info,
	}
}

func (l *Log) Write(p []byte) (n int, err error) {
	zlog.Info(string(p))
	return len(p), nil
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) LogMode(level gLogger.LogLevel) gLogger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *Log) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Info {
		Infof(msg, data...)
	}
}

func (l *Log) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Warn {
		Warnf(msg, data...)
	}
}

func (l *Log) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gLogger.Error {
		Errorf(msg, data...)
	}
}

// Each time an SQL statement is executed, the `trace` method is called.
func (l *Log) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= gLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fileWithLineNum := gUtils.FileWithLineNum()
	logMessage := fmt.Sprintf("%s [%v Rows %d] %s", fileWithLineNum, elapsed, rows, sql)
	if err != nil {
		if logFormat == JsonMode {
			logMessage = fmt.Sprintf("%s  %v [%v Rows %d] %s", fileWithLineNum, err, elapsed, rows, sql)
		} else {
			logMessage = fmt.Sprintf("%s %s %v %s [%v Rows %d] %s",
				ColorRed, fileWithLineNum, err, ColorReset, elapsed, rows, sql)
		}
	}
	switch {
	case l.LogLevel >= gLogger.Info:
		Infof(logMessage)
	case err != nil && l.LogLevel >= gLogger.Error:
		Errorf(logMessage)
	}
}
func InterceptorLogger() logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		switch level {
		case logging.LevelDebug:
			Debugf(msg, fields...)
		case logging.LevelInfo:
			Infof(msg, fields...)
		case logging.LevelWarn:
			Warnf(msg, fields...)
		case logging.LevelError:
			Errorf(msg, fields...)
		}
	})
}
