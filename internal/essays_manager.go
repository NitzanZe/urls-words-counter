package internal

import (
	"context"
	"github.com/NitzanZe/urls-words-counter/internal/helpers"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

// EssaysManager This struct is for managing the essays list
type EssaysManager struct {
	ctx        context.Context
	logger     *zap.SugaredLogger
	urls       *cache.Cache
	wordsCache *cache.Cache
	wordChan   chan string
}

// NewEssaysManager Will init new Essays mananger
func NewEssaysManager(ctx context.Context, logger *zap.SugaredLogger) *EssaysManager {
	return &EssaysManager{
		ctx:    ctx,
		logger: logger,
		urls:   cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// LoadEssaysUrlsFromFileIntoMemory Load all the essays into memory using the helper function
func (em *EssaysManager) LoadEssaysUrlsFromFileIntoMemory(fileName string) error {
	count, duration, err := helpers.ScanFileAndLoadIntoCache(em.logger, fileName, em.urls, false, nil)
	if err != nil {
		return err
	}
	em.logger.Infof("Loaded total of %d urls into memory in %f seconds", count, duration)
	return nil
}

// SetWordsCache will set the words cache into the EssatsManager
func (em *EssaysManager) SetWordsCache(wordsCache *cache.Cache) {
	em.wordsCache = wordsCache
}

// ExtractWordsMapFromUrl Extracting all the words from the url using html tokenizer
func (em *EssaysManager) ExtractWordsMapFromUrl(url string) (map[string]int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	tokenizer := html.NewTokenizer(resp.Body)
	wordsMap := make(map[string]int)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return wordsMap, nil
		case html.TextToken:
			// Process each word in the text
			words := strings.Fields(tokenizer.Token().Data)
			for _, word := range words {
				if helpers.IsValidWord(word) {
					_, ok := em.wordsCache.Get(word)
					if ok {
						wordsMap[strings.ToLower(word)] += 1
					}
				}
			}
		}
	}
}
