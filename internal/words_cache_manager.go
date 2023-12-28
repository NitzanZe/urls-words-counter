package internal

import (
	"context"
	"github.com/NitzanZe/urls-words-counter/internal/helpers"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
)

// WordsManager This struct is for managing the words list
type WordsManager struct {
	ctx    context.Context
	logger *zap.SugaredLogger
	words  *cache.Cache
}

// NewWordsManager Init new words manager
func NewWordsManager(ctx context.Context, logger *zap.SugaredLogger) *WordsManager {
	return &WordsManager{
		ctx:    ctx,
		logger: logger,
		words:  cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// LoadAllWordsFromFileIntoCache Laoding all the words from the list into the memory
func (wc *WordsManager) LoadAllWordsFromFileIntoCache(fileName string) error {
	count, duration, err := helpers.ScanFileAndLoadIntoCache(wc.logger, fileName, wc.words, true, helpers.IsValidWord)
	if err != nil {
		return err
	}
	wc.logger.Infof("Loaded total of %d urls into memory in %f seconds", count, duration)
	return nil
}

// GetWordsCache Return a pointer to the words cache
func (wc *WordsManager) GetWordsCache() *cache.Cache {
	return wc.words
}
