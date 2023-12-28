package main

import (
	"context"
	"fmt"
	"github.com/NitzanZe/urls-words-counter/internal"
	"github.com/NitzanZe/urls-words-counter/internal/config"
	"github.com/NitzanZe/urls-words-counter/internal/models"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("WordsFileFullPath and UrlsFileFullPath wasn't provided. Please provide them by adding them to the run command.\nFor example ./main /home/user/words.txt /home/user/endg-urls")
		return
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	/*
		==================== INITIALIZE LOGGER AND CONFIG ================
	*/
	logger, cfg, err := config.HandleConfig(ctx)
	if err != nil {
		logger.Fatalf("Failed to handle config: %v", err)
	}
	cfg.General.WordsFileFullPath = os.Args[1]
	cfg.General.UrlsFileFullPath = os.Args[2]

	errorChan := make(chan models.Error, 1)

	/*
		==================================================================
		Controller to initialize and control all relevant components
		==================================================================
	*/
	controller := internal.NewController(ctx, cfg, logger, errorChan)
	err = controller.StartComponents()
	if err != nil {
		logger.Fatal(err)
	}
	//Start workers for processing urls in parallel
	controller.StartWorkers()

	//Start processing url's
	controller.StartUrlProcessingAsync()

	outputChan := controller.MergeResultsAsync()

	/*
		==================================================================
		listen to errors and signals from all components
		==================================================================
	*/
	shutDownChannel := make(chan os.Signal, 1)
	signal.Notify(shutDownChannel, syscall.SIGINT, syscall.SIGTERM)
	internal.ListenToSignals(ctx, shutDownChannel, logger, cancelFunc, outputChan, errorChan)

}
