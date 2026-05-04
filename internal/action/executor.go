package action

import (
	"strings"

	"devmemory/internal/core"
)

// DetectType infers EntryType from content using simple heuristics.
func DetectType(content string) core.EntryType {
	trimmed := strings.TrimSpace(content)

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return core.EntryTypeURL
	}

	shellIndicators := []string{"|", "&&", "||", ";", "$", "`"}
	for _, ch := range shellIndicators {
		if strings.Contains(trimmed, ch) {
			return core.EntryTypeCommand
		}
	}

	return core.EntryTypeNote
}

// IsDangerous checks whether a command string matches known dangerous patterns.
func IsDangerous(cmd string) bool {
	lower := strings.ToLower(cmd)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}
	return false
}

var dangerousPatterns = []string{
	"rm -rf", "rm -fr",
	"del /s", "del /q",
	"format ",
	"shutdown",
	"mkfs.",
	":(){ :|:& };:",
	"dd if=",
	"> /dev/sd",
	"chmod -r 777 /",
}
