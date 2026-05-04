package search

import (
	"devmemory/internal/core"
	"sort"
	"strings"
)

// SearchResult wraps an entry with its search score.
type SearchResult struct {
	Entry *core.Entry
	Score float64
}

// Engine performs text search over entries.
type Engine struct{}

// NewEngine creates a new search engine.
func NewEngine() *Engine {
	return &Engine{}
}

// Search scores all entries against the query and returns sorted results.
func (e *Engine) Search(entries []*core.Entry, query string) []*SearchResult {
	terms := splitTerms(query)
	if len(terms) == 0 {
		return nil
	}

	var results []*SearchResult
	for _, entry := range entries {
		score := Score(entry, terms, DefaultWeights)
		if score > 0 {
			results = append(results, &SearchResult{
				Entry: entry,
				Score: score,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

func splitTerms(query string) []string {
	return strings.Fields(strings.ToLower(query))
}
