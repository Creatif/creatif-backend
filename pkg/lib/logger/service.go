package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var packageInfoLogger *zap.SugaredLogger
var packageErrorLogger *zap.SugaredLogger
var packageWarningLogger *zap.SugaredLogger

func wrapLumberjack(level zapcore.Level, fileName string) func(core zapcore.Core) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     10,
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		level,
	)

	return func(core2 zapcore.Core) zapcore.Core {
		return core
	}
}

func buildBaseLogger(logDir string, level zapcore.Level, fileName string) (*zap.SugaredLogger, error) {
	logFile := fmt.Sprintf("%s/%s", logDir, fileName)

	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{logFile}
	cfg.ErrorOutputPaths = []string{"stderr"}
	cfg.DisableStacktrace = false
	cfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := cfg.Build(zap.WrapCore(wrapLumberjack(level, logFile)))

	if err != nil {
		return nil, err
	}

	createdLogger := logger.Sugar()

	err = createdLogger.Sync()

	if err != nil {
		return nil, err
	}

	return createdLogger, nil
}

func buildInfoLogger(logDir string) error {
	infoLogger, err := buildBaseLogger(logDir, zap.InfoLevel, "info.log")
	if err != nil {
		return err
	}

	packageInfoLogger = infoLogger

	return nil
}

func buildErrorLogger(logDir string) error {
	errorLogger, err := buildBaseLogger(logDir, zap.ErrorLevel, "error.log")
	if err != nil {
		return err
	}

	packageErrorLogger = errorLogger

	return nil
}

func buildWarningLogger(logDir string) error {
	warningLogger, err := buildBaseLogger(logDir, zap.WarnLevel, "warn.log")
	if err != nil {
		return err
	}

	packageWarningLogger = warningLogger

	return nil
}

func BuildLoggers(logDir string) error {
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err := os.MkdirAll(logDir, os.ModePerm)

		if err != nil {
			return err
		}
	}

	if err := buildInfoLogger(logDir); err != nil {
		return err
	}

	if err := buildErrorLogger(logDir); err != nil {
		return err
	}

	if err := buildWarningLogger(logDir); err != nil {
		return err
	}

	return nil
}

func Info(msg ...interface{}) {
	packageInfoLogger.Info(msg)
}

func Error(msg ...interface{}) {
	packageErrorLogger.Error(msg)
}

func Warn(msg ...interface{}) {
	packageWarningLogger.Warn(msg)
}
