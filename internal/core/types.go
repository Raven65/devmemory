package core

// EntryType defines the kind of entry.
type EntryType string

const (
	EntryTypeCommand  EntryType = "command"
	EntryTypeURL      EntryType = "url"
	EntryTypeSnippet  EntryType = "snippet"
	EntryTypePrompt   EntryType = "prompt"
	EntryTypeNote     EntryType = "note"
	EntryTypeIssue    EntryType = "issue"
	EntryTypeJournal  EntryType = "journal"
	EntryTypeBusiness EntryType = "business"
	EntryTypeTask     EntryType = "task"
	EntryTypeFile     EntryType = "file"
	EntryTypeFolder   EntryType = "folder"
)

// ValidEntryTypes returns all valid entry types.
func ValidEntryTypes() []EntryType {
	return []EntryType{
		EntryTypeCommand, EntryTypeURL, EntryTypeSnippet, EntryTypePrompt,
		EntryTypeNote, EntryTypeIssue, EntryTypeJournal, EntryTypeBusiness,
		EntryTypeTask, EntryTypeFile, EntryTypeFolder,
	}
}

// IsValid checks whether the EntryType is a known value.
func (t EntryType) IsValid() bool {
	for _, vt := range ValidEntryTypes() {
		if t == vt {
			return true
		}
	}
	return false
}

// String implements fmt.Stringer.
func (t EntryType) String() string {
	return string(t)
}
