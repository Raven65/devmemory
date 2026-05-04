package core

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

// Entry is the unified model for all records in DevMemory.
type Entry struct {
	ID         string     `json:"id"`
	Type       EntryType  `json:"type"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	Summary    string     `json:"summary,omitempty"`
	Project    string     `json:"project,omitempty"`
	Tags       []string   `json:"tags,omitempty"`
	Favorite   bool       `json:"favorite"`
	Dangerous  bool       `json:"dangerous"`
	Archived   bool       `json:"archived"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	UseCount   int        `json:"use_count"`
}

// generateID creates a random 32-char hex string.
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// NewEntry creates a new Entry with sensible defaults.
func NewEntry(entryType EntryType, content string) *Entry {
	now := time.Now()
	if !entryType.IsValid() {
		entryType = EntryTypeNote
	}
	return &Entry{
		ID:        generateID(),
		Type:      entryType,
		Content:   content,
		Tags:      []string{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TitleOrContent returns the title if set, otherwise a truncated first line of content.
func (e *Entry) TitleOrContent() string {
	if e.Title != "" {
		return e.Title
	}
	firstLine := e.Content
	if idx := strings.Index(e.Content, "\n"); idx > 0 {
		firstLine = e.Content[:idx]
	}
	if len(firstLine) > 80 {
		return firstLine[:80] + "..."
	}
	return firstLine
}
