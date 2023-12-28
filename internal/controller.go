package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/NitzanZe/urls-words-counter/internal/config"
	"github.com/NitzanZe/urls-words-counter/internal/helpers"
	"github.com/NitzanZe/urls-words-counter/internal/models"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
	"runtime"
	"sync"
	"time"
)

type Controller struct {
	ctx           context.Context
	logger        *zap.SugaredLogger
	config        *config.Config
	wordsManager  *WordsManager
	essaysManager *EssaysManager
	workersNumber int
	rateLimiter   ratelimit.Limiter
	wg            sync.WaitGroup
	urlChan       chan string
	resultChan    chan map[string]int
	resultMap     map[string]int
	doneChan      chan bool
	errorChan     chan models.Error
}

// NewController Create new controller for the application
func NewController(ctx context.Context, config *config.Config, logger *zap.SugaredLogger, errorChan chan models.Error) *Controller {
	controller := &Controller{
		ctx:           ctx,
		logger:        logger,
		config:        config,
		wordsManager:  NewWordsManager(ctx, logger),
		essaysManager: NewEssaysManager(ctx, logger),
		rateLimiter:   ratelimit.New(config.General.MaximumWorkersRequestsPerSecond),
		workersNumber: runtime.NumCPU() * config.General.WorkersMultiplier, // The number of workers will be determined by the CPU cores. Setting this to more than one will allow more go-routines for each core
		resultMap:     make(map[string]int),
		resultChan:    make(chan map[string]int),
		doneChan:      make(chan bool),
		errorChan:     errorChan,
	}
	controller.urlChan = make(chan string, controller.workersNumber)

	return controller
}

// StartComponents Will load everything that is needed into the memory
func (controller *Controller) StartComponents() error {
	controller.logger.Infof("Going to run words scanner and load valid words into memory")
	err := controller.wordsManager.LoadAllWordsFromFileIntoCache(controller.config.General.WordsFileFullPath)
	if err != nil {
		return err
	}

	controller.logger.Infof("Going to run urls scanner and load valid urls into memory")
	err = controller.essaysManager.LoadEssaysUrlsFromFileIntoMemory(controller.config.General.UrlsFileFullPath)
	if err != nil {
		return err
	}
	controller.essaysManager.SetWordsCache(controller.wordsManager.GetWordsCache())

	return nil
}

// StartWorkers Will start the workers in seperated go routines.
func (controller *Controller) StartWorkers() {
	controller.logger.Infof("Starting %d url workers to process url and count words.", controller.workersNumber)
	for i := 0; i < controller.workersNumber; i++ {
		controller.wg.Add(1)
		go controller.worker()
	}
}

// StartUrlProcessingAsync Start the url list processing. Will signal the workers to stop at the end by closing the channel and wait for them to finish before sending notification
func (controller *Controller) StartUrlProcessingAsync() {
	go func() {

		for url := range controller.essaysManager.urls.Items() {
			controller.urlChan <- url
		}

		// Signal workers to stop
		close(controller.urlChan)

		// Wait for all the workers to finish
		controller.wg.Wait()

		controller.doneChan <- true
	}()
}

// MergeResultsAsync Will merge the results from all the workers
func (controller *Controller) MergeResultsAsync() chan []models.WordCount {
	outputChan := make(chan []models.WordCount)
	var resultsCounter int
	//The next go routine is for printing to screen how many uls were processed so the user will see that the service is still running
	go func() {
		for {
			before := time.Now()
			select {
			case <-controller.ctx.Done():
				return
			case <-time.After(3*time.Second - time.Since(before)):
				controller.logger.Infof("Processed total of %d urls.", resultsCounter)
			}
		}
	}()
	go func() {

		for {
			select {
			case <-controller.ctx.Done():
				return
			case <-controller.doneChan:
				outputChan <- helpers.GetTopWords(controller.resultMap, controller.config.General.GetTopNWords)

			case result := <-controller.resultChan:
				resultsCounter += 1
				helpers.MergeWordCountResults(controller.resultMap, result)
			}

		}
	}()
	return outputChan
}

// worker functionality for each worker goroutine
func (controller *Controller) worker() {
	defer controller.wg.Done()
	for url := range controller.urlChan {
		// Limit the number of requests per second
		controller.rateLimiter.Take()
		wordCount, err := controller.essaysManager.ExtractWordsMapFromUrl(url)
		if err != nil {
			controller.errorChan <- models.Error{
				Err:   errors.New(fmt.Sprintf("Got an error while trying to process url: %s. This url will not be included in the count", url)),
				Level: models.WARN,
			}
			continue
		}
		controller.logger.Debugf("Top words for url %v: %v", url, wordCount)
		controller.resultChan <- wordCount
	}
}
