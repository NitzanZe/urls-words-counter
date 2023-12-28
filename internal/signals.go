package internal

import (
	"context"
	"encoding/json"
	"github.com/NitzanZe/urls-words-counter/internal/models"
	"go.uber.org/zap"
	"os"
	"time"
)

// ListenToSignals Fan in for signals that arrives from different go routines, including errors anc output handling
func ListenToSignals(ctx context.Context, shutDownChannel chan os.Signal, logger *zap.SugaredLogger, cancelFunc context.CancelFunc, outputChan chan []models.WordCount, errChannel <-chan models.Error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-shutDownChannel:
			cancelFunc()
			logger.Fatal("Got shutdown signal. Terminating...")
		case err := <-errChannel:
			if err.Err != nil {
				logger.Warnf("Got an Error in the ErrorChannel. Err = %v, level = %v", err.Err, err.Level)
				if err.Level == models.FATAL {
					logger.Fatal("The error is fatal. Terminating")
				}
			}
		case output := <-outputChan:
			byteArrOutput, err := json.MarshalIndent(output, "", " ")
			if err != nil {
				logger.Warnf("Failed to convert result to json")
			}
			logger.Infof(string(byteArrOutput))
			cancelFunc()
			// A little sleep for graceful shutdown
			time.Sleep(time.Second)
		}
	}
}
