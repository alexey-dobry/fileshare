package zap

import (
	"log"
	"os"
	"path/filepath"

	"github.com/alexey-dobry/fileshare/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Dir        string `yaml:"dir"`
	Debug      bool   `yaml:"is_debug"`
	Production bool   `yaml:"is_production"`
}

type zapLogger struct {
	logger *zap.SugaredLogger
}

func NewLogger(cfg Config) logger.Logger {
	logDirPath := "./logs"

	if cfg.Dir != "" {
		logDirPath = cfg.Dir
	}

	os.MkdirAll(logDirPath, os.ModePerm)

	logFile, err := os.OpenFile(filepath.Join(logDirPath, "main.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Failed to create log file(%s): %s", logFile.Name(), err)
	}

	var zcfg zapcore.EncoderConfig
	if cfg.Production {
		zcfg = zap.NewProductionEncoderConfig()
	} else {
		zcfg = zap.NewDevelopmentEncoderConfig()
	}

	zcfg.EncodeTime = zapcore.ISO8601TimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(zcfg)

	zcfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(zcfg)

	var levelEncoder zapcore.Level
	if cfg.Debug {
		levelEncoder = zapcore.DebugLevel
	} else {
		levelEncoder = zapcore.ErrorLevel
	}
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), levelEncoder),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), levelEncoder),
	}

	core := zapcore.NewTee(cores...)

	return &zapLogger{
		logger: zap.New(core).Sugar(),
	}
}

func (l *zapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l *zapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l *zapLogger) WithFields(args ...interface{}) logger.Logger {
	zl := zapLogger{
		logger: l.logger.With(args...),
	}
	return &zl
}
