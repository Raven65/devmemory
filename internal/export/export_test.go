package export

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"devmemory/internal/core"
)

func TestExportJSON(t *testing.T) {
	entries := []*core.Entry{
		core.NewEntry(core.EntryTypeCommand, "ls -la"),
		core.NewEntry(core.EntryTypeURL, "https://example.com"),
	}

	data, err := ExportJSON(entries)
	if err != nil {
		t.Fatalf("ExportJSON error: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("ExportJSON returned empty data")
	}

	// Verify it's valid JSON
	var parsed []*core.Entry
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("ExportJSON output is not valid JSON: %v", err)
	}

	if len(parsed) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(parsed))
	}
}

func TestImportJSON(t *testing.T) {
	original := []*core.Entry{
		core.NewEntry(core.EntryTypeCommand, "echo hello"),
		core.NewEntry(core.EntryTypeNote, "test note"),
	}

	data, err := ExportJSON(original)
	if err != nil {
		t.Fatalf("ExportJSON error: %v", err)
	}

	imported, err := ImportJSON(data)
	if err != nil {
		t.Fatalf("ImportJSON error: %v", err)
	}

	if len(imported) != len(original) {
		t.Fatalf("expected %d entries, got %d", len(original), len(imported))
	}

	for i, e := range imported {
		if e.ID != original[i].ID {
			t.Errorf("entry %d: expected ID %s, got %s", i, original[i].ID, e.ID)
		}
		if e.Content != original[i].Content {
			t.Errorf("entry %d: expected Content %q, got %q", i, original[i].Content, e.Content)
		}
		if e.Type != original[i].Type {
			t.Errorf("entry %d: expected Type %s, got %s", i, original[i].Type, e.Type)
		}
	}
}

func TestRoundTripWithAllFields(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	original := core.NewEntry(core.EntryTypeCommand, "rm -rf /tmp/test")
	original.Title = "Clean temp"
	original.Project = "sysadmin"
	original.Tags = []string{"linux", "cleanup"}
	original.Favorite = true
	original.Dangerous = true
	original.Archived = false
	original.UseCount = 5
	original.CreatedAt = now
	original.UpdatedAt = now

	data, err := ExportJSON([]*core.Entry{original})
	if err != nil {
		t.Fatalf("ExportJSON error: %v", err)
	}

	imported, err := ImportJSON(data)
	if err != nil {
		t.Fatalf("ImportJSON error: %v", err)
	}

	e := imported[0]
	if e.Title != original.Title {
		t.Errorf("Title: expected %q, got %q", original.Title, e.Title)
	}
	if e.Project != original.Project {
		t.Errorf("Project: expected %q, got %q", original.Project, e.Project)
	}
	if len(e.Tags) != len(original.Tags) {
		t.Errorf("Tags: expected %v, got %v", original.Tags, e.Tags)
	}
	if e.Favorite != original.Favorite {
		t.Errorf("Favorite: expected %v, got %v", original.Favorite, e.Favorite)
	}
	if e.Dangerous != original.Dangerous {
		t.Errorf("Dangerous: expected %v, got %v", original.Dangerous, e.Dangerous)
	}
	if e.UseCount != original.UseCount {
		t.Errorf("UseCount: expected %d, got %d", original.UseCount, e.UseCount)
	}
}

func TestImportJSONInvalidData(t *testing.T) {
	_, err := ImportJSON([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestExportDailyMarkdown(t *testing.T) {
	date := time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC)

	entries := []*core.Entry{
		{
			ID:        "abc123",
			Type:      core.EntryTypeCommand,
			Content:   "ls -la",
			CreatedAt: time.Date(2026, 5, 4, 10, 30, 0, 0, time.UTC),
		},
		{
			ID:        "def456",
			Type:      core.EntryTypeURL,
			Title:     "Go Docs",
			Content:   "https://go.dev/doc/",
			Project:   "learning",
			Tags:      []string{"go", "docs"},
			CreatedAt: time.Date(2026, 5, 4, 14, 0, 0, 0, time.UTC),
		},
	}

	md, err := ExportDailyMarkdown(entries, date)
	if err != nil {
		t.Fatalf("ExportDailyMarkdown error: %v", err)
	}

	if !strings.Contains(md, "2026-05-04") {
		t.Error("markdown should contain the date")
	}
	if !strings.Contains(md, "10:30") {
		t.Error("markdown should contain time of first entry")
	}
	if !strings.Contains(md, "ls -la") {
		t.Error("markdown should contain entry content")
	}
	if !strings.Contains(md, "Go Docs") {
		t.Error("markdown should contain entry title")
	}
	if !strings.Contains(md, "`learning`") {
		t.Error("markdown should contain project")
	}
	if !strings.Contains(md, "`go`") {
		t.Error("markdown should contain tags")
	}
}

func TestExportDailyMarkdownEmpty(t *testing.T) {
	date := time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC)
	md, err := ExportDailyMarkdown(nil, date)
	if err != nil {
		t.Fatalf("ExportDailyMarkdown error: %v", err)
	}
	if !strings.Contains(md, "no entries") {
		t.Error("markdown should indicate no entries")
	}
}
