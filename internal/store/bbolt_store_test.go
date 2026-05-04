package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"devmemory/internal/core"
)

func newTestStore(t *testing.T) *BBoltStore {
	t.Helper()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")
	s := NewBBoltStore(dbPath)
	if err := s.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestBBoltStore_CreateAndGet(t *testing.T) {
	s := newTestStore(t)

	entry := core.NewEntry(core.EntryTypeCommand, "ls -la")
	entry.Tags = []string{"linux", "cli"}
	entry.Project = "test"

	if err := s.Create(entry); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := s.Get(entry.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}

	if got.Content != "ls -la" {
		t.Errorf("Content: got %q, want %q", got.Content, "ls -la")
	}
	if got.Type != core.EntryTypeCommand {
		t.Errorf("Type: got %q, want %q", got.Type, core.EntryTypeCommand)
	}
	if got.Project != "test" {
		t.Errorf("Project: got %q, want %q", got.Project, "test")
	}
	if len(got.Tags) != 2 {
		t.Errorf("Tags: got %d, want 2", len(got.Tags))
	}
}

func TestBBoltStore_GetNotFound(t *testing.T) {
	s := newTestStore(t)

	_, err := s.Get("nonexistent")
	if err == nil {
		t.Error("Get nonexistent: expected error, got nil")
	}
}

func TestBBoltStore_Update(t *testing.T) {
	s := newTestStore(t)

	entry := core.NewEntry(core.EntryTypeCommand, "ls -la")
	s.Create(entry)

	entry.Content = "ls -la --color"
	entry.Tags = []string{"linux", "cli", "color"}
	if err := s.Update(entry); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := s.Get(entry.ID)
	if err != nil {
		t.Fatalf("Get after Update: %v", err)
	}
	if got.Content != "ls -la --color" {
		t.Errorf("Content after Update: got %q, want %q", got.Content, "ls -la --color")
	}
	if len(got.Tags) != 3 {
		t.Errorf("Tags after Update: got %d, want 3", len(got.Tags))
	}
}

func TestBBoltStore_UpdateNotFound(t *testing.T) {
	s := newTestStore(t)

	entry := core.NewEntry(core.EntryTypeCommand, "ls")
	entry.ID = "nonexistent"
	if err := s.Update(entry); err == nil {
		t.Error("Update nonexistent: expected error, got nil")
	}
}

func TestBBoltStore_Delete(t *testing.T) {
	s := newTestStore(t)

	entry := core.NewEntry(core.EntryTypeCommand, "ls -la")
	s.Create(entry)

	if err := s.Delete(entry.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := s.Get(entry.ID)
	if err == nil {
		t.Error("Get after Delete: expected error, got nil")
	}
}

func TestBBoltStore_DeleteNotFound(t *testing.T) {
	s := newTestStore(t)

	if err := s.Delete("nonexistent"); err == nil {
		t.Error("Delete nonexistent: expected error, got nil")
	}
}

func TestBBoltStore_List(t *testing.T) {
	s := newTestStore(t)

	s.Create(core.NewEntry(core.EntryTypeCommand, "ls -la"))
	s.Create(core.NewEntry(core.EntryTypeURL, "https://example.com"))
	s.Create(core.NewEntry(core.EntryTypeCommand, "grep -r pattern"))

	// All entries
	entries, err := s.List(ListOptions{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("List all: got %d, want 3", len(entries))
	}

	// Filter by type
	entries, err = s.List(ListOptions{Type: core.EntryTypeCommand})
	if err != nil {
		t.Fatalf("List by type: %v", err)
	}
	if len(entries) != 2 {
		t.Errorf("List by type: got %d, want 2", len(entries))
	}

	// Filter by non-matching type
	entries, err = s.List(ListOptions{Type: core.EntryTypeNote})
	if err != nil {
		t.Fatalf("List by type (no match): %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("List by type (no match): got %d, want 0", len(entries))
	}
}

func TestBBoltStore_ListFilterByTag(t *testing.T) {
	s := newTestStore(t)

	e1 := core.NewEntry(core.EntryTypeCommand, "ls -la")
	e1.Tags = []string{"linux", "cli"}
	s.Create(e1)

	e2 := core.NewEntry(core.EntryTypeCommand, "grep pattern")
	e2.Tags = []string{"linux", "search"}
	s.Create(e2)

	entries, err := s.List(ListOptions{Tag: "cli"})
	if err != nil {
		t.Fatalf("List by tag: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("List by tag: got %d, want 1", len(entries))
	}
}

func TestBBoltStore_ListFilterByProject(t *testing.T) {
	s := newTestStore(t)

	e1 := core.NewEntry(core.EntryTypeNote, "note 1")
	e1.Project = "alpha"
	s.Create(e1)

	e2 := core.NewEntry(core.EntryTypeNote, "note 2")
	e2.Project = "beta"
	s.Create(e2)

	entries, err := s.List(ListOptions{Project: "alpha"})
	if err != nil {
		t.Fatalf("List by project: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("List by project: got %d, want 1", len(entries))
	}
}

func TestBBoltStore_GetByDate(t *testing.T) {
	s := newTestStore(t)

	e1 := core.NewEntry(core.EntryTypeCommand, "today command")
	s.Create(e1)

	now := time.Now()
	entries, err := s.GetByDate(now.Year(), int(now.Month()), now.Day())
	if err != nil {
		t.Fatalf("GetByDate: %v", err)
	}
	if len(entries) != 1 {
		t.Errorf("GetByDate today: got %d, want 1", len(entries))
	}

	// Yesterday should have no entries
	yesterday := now.Add(-24 * time.Hour)
	entries, err = s.GetByDate(yesterday.Year(), int(yesterday.Month()), yesterday.Day())
	if err != nil {
		t.Fatalf("GetByDate yesterday: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("GetByDate yesterday: got %d, want 0", len(entries))
	}
}

func TestBBoltStore_OpenCreatesDir(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "subdir", "test.db")

	s := NewBBoltStore(dbPath)
	if err := s.Open(); err != nil {
		t.Fatalf("Open with new dir: %v", err)
	}
	s.Close()

	if _, err := os.Stat(filepath.Join(dir, "subdir")); os.IsNotExist(err) {
		t.Error("Open did not create subdirectory")
	}
}
