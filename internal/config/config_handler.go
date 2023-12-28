package config

import (
	"context"
	zapLogger "github.com/NitzanZe/urls-words-counter/pkg/logger"
	"github.com/codingconcepts/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var ServiceName = "urls-top-words-counter"

func setConfigEnv(config *Config) error {
	fields := []interface{}{&config.Logger, &config.General}

	for _, field := range fields {
		err := env.Set(field)
		if err != nil {
			return err
		}
	}
	return nil
}

func HandleConfig(ctx context.Context) (*zap.SugaredLogger, *Config, error) {
	config := &Config{}
	err := setConfigEnv(config)
	if err != nil {
		return nil, nil, err
	}

	err = os.MkdirAll(config.Logger.LogPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	today := time.Now().Format("2006_01_02")
	logPattern := "/" + today + ".log"
	file, err := os.Create(config.Logger.LogPath + logPattern)
	if err != nil {
		return nil, nil, err
	}
	var outputPath []string
	if config.Logger.LogEnableStdOutput {
		outputPath = append(outputPath, "stdout")
	}
	if config.Logger.LogEnableFileOutput {
		zapcore.AddSync(file)
		outputPath = append(outputPath, config.Logger.LogPath+logPattern)
	}
	log, err := zapLogger.NewLogger(ServiceName, config.Logger.LogLevel, config.Logger.LogFormatting, outputPath...)
	if err != nil {
		return nil, nil, err
	}
	go func() {
		<-ctx.Done()
		log.Infof("Flushing logger and closing log file")
		log.Sync()
		file.Close()
	}()
	return log, config, nil
}
