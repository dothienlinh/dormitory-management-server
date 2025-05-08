package pkg

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *zap.Logger

// LogConfig contains the configuration for the logger
type LogConfig struct {
	Environment    string // "development" or "production"
	LogLevel       string // "debug", "info", "warn", "error"
	LogPath        string // log path, empty if only log to console
	MaxSize        int    // MB
	MaxBackups     int    // max number of backup files
	MaxAge         int    // max number of days to keep the log file
	Compress       bool   // compress log file
	JsonFormat     bool   // log format as JSON
	ShowCaller     bool   // show caller in log
	ShowStacktrace bool   // show stacktrace for log level >= error
}

// NewDefaultLogConfig returns the default configuration for the logger
func NewDefaultLogConfig() LogConfig {
	return LogConfig{
		Environment:    "development",
		LogLevel:       "info",
		LogPath:        "",
		MaxSize:        100,
		MaxBackups:     5,
		MaxAge:         30,
		Compress:       true,
		JsonFormat:     false,
		ShowCaller:     true,
		ShowStacktrace: true,
	}
}

// NewLogger creates and returns a logger with custom configuration
func NewLogger(config ...LogConfig) *zap.Logger {
	cfg := NewDefaultLogConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Setup Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Better configuration for development environment
	if cfg.Environment == "development" && !cfg.JsonFormat {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
	}

	// Determine encoder (JSON or console)
	var encoder zapcore.Encoder
	if cfg.JsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// Determine level
	var level zapcore.Level
	switch cfg.LogLevel {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Setup output (console and/or file)
	var cores []zapcore.Core

	// Always log to console
	consoleWriter := zapcore.Lock(os.Stdout)
	cores = append(cores, zapcore.NewCore(encoder, consoleWriter, level))

	// Add log to file if configured
	if cfg.LogPath != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.LogPath,
			MaxSize:    cfg.MaxSize,    // MB
			MaxBackups: cfg.MaxBackups, // number of backup files
			MaxAge:     cfg.MaxAge,     // days
			Compress:   cfg.Compress,   // compress file
		})
		cores = append(cores, zapcore.NewCore(encoder, fileWriter, level))
	}

	// Combine all cores
	core := zapcore.NewTee(cores...)

	// Create logger
	var options []zap.Option
	if cfg.ShowCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.ShowStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	logger := zap.New(core, options...)

	// Save logger to global variable to be used anywhere
	globalLogger = logger

	return logger
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	if globalLogger == nil {
		globalLogger = NewLogger()
	}
	return globalLogger
}

// WithField adds a field to the logger
func WithField(key string, value interface{}) *zap.Logger {
	return GetLogger().With(zap.Any(key, value))
}

// WithError adds an error to the logger
func WithError(err error) *zap.Logger {
	return GetLogger().With(zap.Error(err))
}

// WithRequest adds request information to the logger
func WithRequest(method, path, ip string) *zap.Logger {
	return GetLogger().With(
		zap.String("method", method),
		zap.String("path", path),
		zap.String("ip", ip),
	)
}

// InitProductionLogger initializes the logger for the production environment
func InitProductionLogger(logPath string) *zap.Logger {
	return NewLogger(LogConfig{
		Environment: "production",
		LogLevel:    "info",
		LogPath:     logPath,
		JsonFormat:  true,
		ShowCaller:  true,
	})
}

// InitDevelopmentLogger initializes the logger for the development environment
func InitDevelopmentLogger() *zap.Logger {
	return NewLogger(LogConfig{
		Environment: "development",
		LogLevel:    "debug",
		JsonFormat:  false,
		ShowCaller:  true,
	})
}

// ConfigureLogger configures the logger from the environment variable
func ConfigureLogger() *zap.Logger {
	env := os.Getenv("APP_ENV")
	logPath := os.Getenv("LOG_PATH")

	if env == "production" {
		if logPath == "" {
			logPath = fmt.Sprintf("./logs/app_%s.log", time.Now().Format("2006-01-02"))
		}
		return InitProductionLogger(logPath)
	}

	return InitDevelopmentLogger()
}
