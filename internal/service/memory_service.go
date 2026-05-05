package service

import (
	"fmt"
	"time"

	"devmemory/internal/action"
	"devmemory/internal/core"
	"devmemory/internal/export"
	"devmemory/internal/search"
	"devmemory/internal/store"
)

// MemoryService provides a unified business layer for CLI and Fyne UI.
type MemoryService struct {
	store    store.Store
	searcher *search.Engine
}

// NewMemoryService creates a new service with the given store.
func NewMemoryService(s store.Store) *MemoryService {
	return &MemoryService{
		store:    s,
		searcher: search.NewEngine(),
	}
}

// CreateEntryInput holds the fields for creating a new entry.
type CreateEntryInput struct {
	Content string
	Type    core.EntryType // empty = auto-detect
	Title   string
	Project string
	Tags    []string
}

// CreateEntry creates a new entry with auto type detection and danger flagging.
func (s *MemoryService) CreateEntry(input CreateEntryInput) (*core.Entry, error) {
	if input.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	entryType := input.Type
	if entryType == "" {
		entryType = action.DetectType(input.Content)
	}

	entry := core.NewEntry(entryType, input.Content)
	if input.Title != "" {
		entry.Title = input.Title
	}
	if input.Project != "" {
		entry.Project = input.Project
	}
	if len(input.Tags) > 0 {
		entry.Tags = input.Tags
	}

	if entry.Type == core.EntryTypeCommand && action.IsDangerous(entry.Content) {
		entry.Dangerous = true
	}

	if err := s.store.Create(entry); err != nil {
		return nil, fmt.Errorf("creating entry: %w", err)
	}

	return entry, nil
}

// GetEntry retrieves an entry by ID (supports short ID prefix).
func (s *MemoryService) GetEntry(id string) (*core.Entry, error) {
	fullID, err := s.store.ResolveID(id)
	if err != nil {
		return nil, fmt.Errorf("resolving ID: %w", err)
	}
	entry, err := s.store.Get(fullID)
	if err != nil {
		return nil, fmt.Errorf("getting entry: %w", err)
	}
	return entry, nil
}

// ListEntries returns entries matching the given filter options.
func (s *MemoryService) ListEntries(opts store.ListOptions) ([]*core.Entry, error) {
	return s.store.List(opts)
}

// SearchEntries searches entries by query with optional type filter.
func (s *MemoryService) SearchEntries(query string, entryType core.EntryType) ([]*search.SearchResult, error) {
	opts := store.ListOptions{Type: entryType}
	entries, err := s.store.List(opts)
	if err != nil {
		return nil, fmt.Errorf("listing entries for search: %w", err)
	}
	results := s.searcher.Search(entries, query)
	return results, nil
}

// GetToday returns entries created today.
func (s *MemoryService) GetToday() ([]*core.Entry, error) {
	now := time.Now()
	return s.store.GetByDate(now.Year(), int(now.Month()), now.Day())
}

// CopyEntry marks an entry as used (UseCount++, LastUsedAt update).
// Actual clipboard copy is handled by the UI layer.
func (s *MemoryService) CopyEntry(id string) (*core.Entry, error) {
	entry, err := s.GetEntry(id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	entry.LastUsedAt = &now
	entry.UseCount++
	entry.UpdatedAt = now
	if err := s.store.Update(entry); err != nil {
		return nil, fmt.Errorf("updating entry after copy: %w", err)
	}
	return entry, nil
}

// UpdateEntryFields holds optional fields to update on an entry.
type UpdateEntryFields struct {
	Title    *string
	Content  *string
	Type     *core.EntryType
	Project  *string
	Tags     *[]string
	Favorite *bool
	Archived *bool
}

// UpdateEntry updates specified fields on an entry.
func (s *MemoryService) UpdateEntry(id string, fields UpdateEntryFields) (*core.Entry, error) {
	entry, err := s.GetEntry(id)
	if err != nil {
		return nil, err
	}

	if fields.Title != nil {
		entry.Title = *fields.Title
	}
	if fields.Content != nil {
		entry.Content = *fields.Content
	}
	if fields.Type != nil {
		entry.Type = *fields.Type
	}
	if fields.Project != nil {
		entry.Project = *fields.Project
	}
	if fields.Tags != nil {
		entry.Tags = *fields.Tags
	}
	if fields.Favorite != nil {
		entry.Favorite = *fields.Favorite
	}
	if fields.Archived != nil {
		entry.Archived = *fields.Archived
	}

	entry.UpdatedAt = time.Now()

	if err := s.store.Update(entry); err != nil {
		return nil, fmt.Errorf("updating entry: %w", err)
	}
	return entry, nil
}

// ToggleFavorite toggles the favorite flag on an entry.
func (s *MemoryService) ToggleFavorite(id string) (*core.Entry, error) {
	entry, err := s.GetEntry(id)
	if err != nil {
		return nil, err
	}
	entry.Favorite = !entry.Favorite
	entry.UpdatedAt = time.Now()
	if err := s.store.Update(entry); err != nil {
		return nil, fmt.Errorf("updating entry: %w", err)
	}
	return entry, nil
}

// ToggleArchive toggles the archived flag on an entry.
func (s *MemoryService) ToggleArchive(id string) (*core.Entry, error) {
	entry, err := s.GetEntry(id)
	if err != nil {
		return nil, err
	}
	entry.Archived = !entry.Archived
	entry.UpdatedAt = time.Now()
	if err := s.store.Update(entry); err != nil {
		return nil, fmt.Errorf("updating entry: %w", err)
	}
	return entry, nil
}

// DeleteEntry deletes an entry by ID.
func (s *MemoryService) DeleteEntry(id string) error {
	fullID, err := s.store.ResolveID(id)
	if err != nil {
		return fmt.Errorf("resolving ID: %w", err)
	}
	return s.store.Delete(fullID)
}

// ExportTodayMarkdown generates a Markdown report of today's entries.
func (s *MemoryService) ExportTodayMarkdown() (string, error) {
	entries, err := s.GetToday()
	if err != nil {
		return "", err
	}
	return export.ExportDailyMarkdown(entries, time.Now())
}

// ExportJSON exports all entries as JSON.
func (s *MemoryService) ExportJSON() ([]byte, error) {
	entries, err := s.store.List(store.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing entries: %w", err)
	}
	return export.ExportJSON(entries)
}

// ImportJSON imports entries from JSON data. Returns the number of imported entries.
func (s *MemoryService) ImportJSON(data []byte) (int, error) {
	entries, err := export.ImportJSON(data)
	if err != nil {
		return 0, fmt.Errorf("parsing JSON: %w", err)
	}

	imported := 0
	for _, entry := range entries {
		if err := s.store.Create(entry); err != nil {
			// ID already exists, try update
			if err := s.store.Update(entry); err != nil {
				continue
			}
		}
		imported++
	}
	return imported, nil
}

// Count returns the total number of entries.
func (s *MemoryService) Count() (int, error) {
	entries, err := s.store.List(store.ListOptions{})
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}
