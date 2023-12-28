package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerEncoding string

const (
	Json    LoggerEncoding = "json"
	Console LoggerEncoding = "console"
)

// NewLogger constructs a Sugared Logger that writes to stdout and
// provides human-readable timestamps.
func NewLogger(service string, level string, encoding string, outputPaths ...string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()
	config.Encoding = encoding
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if TranslateStringToEncoding(encoding) == Console {
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.DisableStacktrace = true
	config.Level = zap.NewAtomicLevelAt(TranslateStringLevelToZapLevel(level))
	config.EncoderConfig.ConsoleSeparator = " "
	config.InitialFields = map[string]any{
		"service": service,
	}

	config.OutputPaths = []string{"stdout"}
	if outputPaths != nil {
		config.OutputPaths = outputPaths
	}

	log, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}

func TranslateStringToEncoding(encoding string) LoggerEncoding {
	switch encoding {
	case "json":
		return Json
	case "console":
		return Console
	default:
		return Json
	}
}
func TranslateStringLevelToZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}
