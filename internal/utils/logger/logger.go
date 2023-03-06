package logger

import (
	"context"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/KarolosLykos/sportsnews/config"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, err error, args ...interface{})
	Warnf(ctx context.Context, err error, format string, args ...interface{})
	Error(ctx context.Context, err error, args ...interface{})
	Errorf(ctx context.Context, err error, format string, args ...interface{})
	Panic(ctx context.Context, err error, args ...interface{})
	Panicf(ctx context.Context, err error, format string, args ...interface{})
	Fatal(ctx context.Context, err error, args ...interface{})
	Fatalf(ctx context.Context, err error, format string, args ...interface{})
}

type exportedLogger struct {
	cfg    *config.Config
	logger *logrus.Logger
}

// Init creates a new Logger and initializes it.
func Init(cfg *config.Config) Logger {
	log := &logrus.Logger{
		Out:          os.Stderr,
		Hooks:        make(logrus.LevelHooks),
		ReportCaller: false,
		ExitFunc:     os.Exit,
	}

	log.SetLevel(setLevel(cfg.Logger.LogLevel))
	log.SetFormatter(setFormatter(cfg.Logger.Format))

	return New(cfg, log)
}

// New creates a new Logger.
func New(cfg *config.Config, l *logrus.Logger) Logger {
	return &exportedLogger{cfg: cfg, logger: l}
}

// setFormatter returns the format  of the Logger.
func setFormatter(format string) logrus.Formatter {
	if strings.ToLower(format) == "json" {
		return &logrus.JSONFormatter{}
	}

	return &logrus.TextFormatter{}
}

// setLevel returns the level of the Logger.
func setLevel(lvl string) logrus.Level {
	level, err := logrus.ParseLevel(lvl)
	if err != nil {
		return logrus.InfoLevel
	}

	return level
}

func (l *exportedLogger) Debug(ctx context.Context, args ...interface{}) {
	le := l.parseArgs(ctx, nil, args...)
	le.Debug(args...)
}

func (l *exportedLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	le := l.parseArgs(ctx, nil, args...)
	le.Debugf(format, args...)
}

func (l *exportedLogger) Info(ctx context.Context, args ...interface{}) {
	le := l.parseArgs(ctx, nil, args...)
	le.Info(args...)
}

func (l *exportedLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	le := l.parseArgs(ctx, nil, args...)
	le.Infof(format, args...)
}

func (l *exportedLogger) Warn(ctx context.Context, err error, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Warn(args...)
}

func (l *exportedLogger) Warnf(ctx context.Context, err error, format string, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Warnf(format, args...)
}

func (l *exportedLogger) Error(ctx context.Context, err error, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Error(args...)
}

func (l *exportedLogger) Errorf(ctx context.Context, err error, format string, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Errorf(format, args...)
}

func (l *exportedLogger) Panic(ctx context.Context, err error, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Panic(args...)
}

func (l *exportedLogger) Panicf(ctx context.Context, err error, format string, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Panicf(format, args...)
}

func (l *exportedLogger) Fatal(ctx context.Context, err error, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Fatal(args...)
}

func (l *exportedLogger) Fatalf(ctx context.Context, err error, format string, args ...interface{}) {
	le := l.parseArgs(ctx, err, args...)
	le.Fatalf(format, args...)
}

// parseArgs parses the arguments and adds err field on the logger.
func (l *exportedLogger) parseArgs(_ context.Context, err error, _ ...interface{}) *logrus.Entry {
	e := logrus.NewEntry(l.logger)

	if err != nil {
		e = e.WithError(err)
	}

	return e
}
