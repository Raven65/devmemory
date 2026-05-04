package search

import (
	"math"
	"strings"
	"time"

	"devmemory/internal/core"
)

// FieldWeights controls the scoring weights for each entry field.
type FieldWeights struct {
	Title    float64
	Tag      float64
	Project  float64
	Summary  float64
	Content  float64
	Favorite float64
}

// DefaultWeights are the standard scoring weights.
var DefaultWeights = FieldWeights{
	Title:    10.0,
	Tag:      7.0,
	Project:  5.0,
	Summary:  3.0,
	Content:  1.0,
	Favorite: 2.0,
}

// Score calculates a relevance score for an entry given search terms and weights.
func Score(entry *core.Entry, terms []string, weights FieldWeights) float64 {
	score := 0.0

	for _, term := range terms {
		lower := strings.ToLower(term)

		if containsLower(entry.Title, lower) {
			score += weights.Title
		}

		for _, tag := range entry.Tags {
			if containsLower(tag, lower) {
				score += weights.Tag
				break
			}
		}

		if containsLower(entry.Project, lower) {
			score += weights.Project
		}

		if containsLower(entry.Summary, lower) {
			score += weights.Summary
		}

		if containsLower(entry.Content, lower) {
			score += weights.Content
		}
	}

	// Favorite bonus
	if entry.Favorite {
		score += weights.Favorite
	}

	// UseCount bonus (capped at 3.0)
	score += math.Min(float64(entry.UseCount)*0.1, 3.0)

	// Recency bonus based on LastUsedAt (decays over ~3 weeks)
	if entry.LastUsedAt != nil {
		daysSince := time.Since(*entry.LastUsedAt).Hours() / 24
		recency := math.Max(0, 3.0-daysSince*3.0/21.0)
		score += recency
	}

	return score
}

func containsLower(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), substr)
}
