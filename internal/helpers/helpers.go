package helpers

import (
	"bufio"
	"github.com/NitzanZe/urls-words-counter/internal/models"
	"github.com/patrickmn/go-cache"
	"go.uber.org/zap"
	"os"
	"sort"
	"time"
	"unicode"
)

/*
IsValidWord Helper function to check if the word is valid by the next rules:
1. Contain at least 3 characters.
2. Contain only alphabetic characters.
*/
func IsValidWord(word string) bool {
	if len(word) < 3 {
		return false
	}

	for _, char := range word {
		if !unicode.IsLetter(char) {
			return false
		}
	}

	return true
}

// ScanFileAndLoadIntoCache Will scan all the words/urls from a file and will perform validation if necessary (validationFunc must be provided)
func ScanFileAndLoadIntoCache(logger *zap.SugaredLogger, fileName string, cachePointer *cache.Cache, shouldValidate bool, validationFunc func(string) bool) (int, float64, error) {
	startTime := time.Now()
	file, err := os.Open(fileName)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		scan := scanner.Text()
		// If validation is required, only valid scan checked with validation function will be added into memory
		if shouldValidate {
			if validationFunc != nil && validationFunc(scan) {
				err = cachePointer.Add(scan, 1, cache.NoExpiration)
				if err != nil {
					logger.Errorf("Failed to add %s to into memory.", scan)
				}
			}
			continue
		}
		err = cachePointer.Add(scan, 1, cache.NoExpiration)
		if err != nil {
			logger.Errorf("Failed to add %s to into memory.", scan)
		}
	}
	return cachePointer.ItemCount(), time.Since(startTime).Seconds(), scanner.Err()
}

// GetTopWords Get top n words from map and return in WorCount struct format
func GetTopWords(words map[string]int, n int) []models.WordCount {
	var wordCounts []models.WordCount
	for word, count := range words {
		wordCounts = append(wordCounts, models.WordCount{Word: word, Count: count})
	}
	// Sort word counts in descending order
	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].Count > wordCounts[j].Count
	})
	if len(wordCounts) < n {
		return wordCounts[:len(wordCounts)]
	}
	return wordCounts[:n]
}

// MergeWordCountResults Merging one map to the other
func MergeWordCountResults(target map[string]int, source map[string]int) {
	for key, value := range source {
		target[key] += value
	}
}
