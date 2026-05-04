package search

import (
	"testing"

	"devmemory/internal/core"
)

func newTestEntry(entryType core.EntryType, content, title, project string, tags []string) *core.Entry {
	e := core.NewEntry(entryType, content)
	e.Title = title
	e.Project = project
	if tags != nil {
		e.Tags = tags
	}
	return e
}

func TestSearch_BasicSingleTerm(t *testing.T) {
	entries := []*core.Entry{
		newTestEntry(core.EntryTypeCommand, "ss -tinp | grep ESTAB", "Check TCP connections", "network", []string{"linux", "tcp"}),
		newTestEntry(core.EntryTypeURL, "https://docs.kernel.org/networking/", "Kernel networking docs", "kernel", []string{"linux", "kernel"}),
		newTestEntry(core.EntryTypeNote, "NUMA binding affects RDMA performance", "", "rdma", []string{"numa", "rdma"}),
	}

	engine := NewEngine()

	results := engine.Search(entries, "tcp")
	if len(results) == 0 {
		t.Fatal("Search 'tcp': expected results, got 0")
	}
	// TCP command should score highest (title + tag match)
	if results[0].Entry.Content != "ss -tinp | grep ESTAB" {
		t.Errorf("Search 'tcp' top result: got content %q", results[0].Entry.Content)
	}
}

func TestSearch_MultiTerm(t *testing.T) {
	entries := []*core.Entry{
		newTestEntry(core.EntryTypeNote, "NUMA binding for RDMA", "", "rdma", []string{"numa"}),
		newTestEntry(core.EntryTypeNote, "RDMA benchmark results", "", "rdma", []string{"benchmark"}),
		newTestEntry(core.EntryTypeNote, "TCP tuning guide", "", "network", []string{"tcp"}),
	}

	engine := NewEngine()

	results := engine.Search(entries, "numa rdma")
	if len(results) == 0 {
		t.Fatal("Search 'numa rdma': expected results, got 0")
	}
	// First entry matches both terms, should be highest
	if results[0].Entry.Content != "NUMA binding for RDMA" {
		t.Errorf("Search 'numa rdma' top result: got %q", results[0].Entry.Content)
	}
}

func TestSearch_NoResults(t *testing.T) {
	entries := []*core.Entry{
		newTestEntry(core.EntryTypeCommand, "ls -la", "", "", nil),
	}

	engine := NewEngine()

	results := engine.Search(entries, "zzzznonexistent")
	if len(results) != 0 {
		t.Errorf("Search 'zzzznonexistent': expected 0 results, got %d", len(results))
	}
}

func TestSearch_EmptyQuery(t *testing.T) {
	entries := []*core.Entry{
		newTestEntry(core.EntryTypeCommand, "ls -la", "", "", nil),
	}

	engine := NewEngine()

	results := engine.Search(entries, "")
	if results != nil {
		t.Errorf("Search '': expected nil, got %d results", len(results))
	}
}

func TestSearch_ProjectMatch(t *testing.T) {
	entries := []*core.Entry{
		newTestEntry(core.EntryTypeNote, "config file location", "", "rdma", nil),
		newTestEntry(core.EntryTypeNote, "config file location", "", "kernel", nil),
	}

	engine := NewEngine()

	results := engine.Search(entries, "rdma")
	if len(results) != 1 {
		t.Fatalf("Search 'rdma': expected 1 result, got %d", len(results))
	}
	if results[0].Entry.Project != "rdma" {
		t.Errorf("Search 'rdma' result project: got %q", results[0].Entry.Project)
	}
}

func TestSearch_ScoreOrdering(t *testing.T) {
	entries := []*core.Entry{
		newTestEntry(core.EntryTypeNote, "just mentioning tcp in content", "", "", nil),
		newTestEntry(core.EntryTypeCommand, "tcpdump -i eth0", "TCP dump tool", "network", []string{"tcp", "debug"}),
	}

	engine := NewEngine()

	results := engine.Search(entries, "tcp")
	if len(results) != 2 {
		t.Fatalf("Search 'tcp': expected 2 results, got %d", len(results))
	}
	// Second entry has title + tag + project matches, should score higher
	if results[0].Entry.Title != "TCP dump tool" {
		t.Errorf("Top result should be 'TCP dump tool', got %q", results[0].Entry.Title)
	}
}

func TestDetectType_DoesNotAffectSearch(t *testing.T) {
	// Ensure search works with all entry types
	entries := []*core.Entry{
		newTestEntry(core.EntryTypeCommand, "ls", "", "", nil),
		newTestEntry(core.EntryTypeURL, "https://example.com", "", "", nil),
		newTestEntry(core.EntryTypeNote, "a note about ls command", "", "", nil),
		newTestEntry(core.EntryTypeSnippet, "const x = 1", "", "", nil),
	}

	engine := NewEngine()

	results := engine.Search(entries, "ls")
	// Should match the command and the note
	if len(results) != 2 {
		t.Errorf("Search 'ls': expected 2 results, got %d", len(results))
	}
}
