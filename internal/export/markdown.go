package export

import (
	"fmt"
	"strings"
	"time"

	"devmemory/internal/core"
)

// ExportDailyMarkdown generates a Markdown-formatted timeline of entries for a given date.
func ExportDailyMarkdown(entries []*core.Entry, date time.Time) (string, error) {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# DevMemory Daily Report — %s\n\n", date.Format("2006-01-02")))
	b.WriteString(fmt.Sprintf("**%d entries**\n\n", len(entries)))

	if len(entries) == 0 {
		b.WriteString("*(no entries recorded)*\n")
		return b.String(), nil
	}

	b.WriteString("---\n\n")

	for _, e := range entries {
		timeStr := e.CreatedAt.Format("15:04")
		b.WriteString(fmt.Sprintf("### %s [%s]\n\n", timeStr, e.Type))

		if e.Title != "" {
			b.WriteString(fmt.Sprintf("**%s**\n\n", e.Title))
		}

		// Indent content lines
		lines := strings.Split(e.Content, "\n")
		for _, line := range lines {
			b.WriteString(fmt.Sprintf("> %s\n", line))
		}
		b.WriteString("\n")

		if e.Project != "" {
			b.WriteString(fmt.Sprintf("- Project: `%s`\n", e.Project))
		}
		if len(e.Tags) > 0 {
			b.WriteString(fmt.Sprintf("- Tags: %s\n", formatTagsMarkdown(e.Tags)))
		}
		if e.Dangerous {
			b.WriteString("- **DANGEROUS**\n")
		}
		if e.Favorite {
			b.WriteString("- Favorite\n")
		}

		b.WriteString("\n---\n\n")
	}

	return b.String(), nil
}

func formatTagsMarkdown(tags []string) string {
	parts := make([]string, len(tags))
	for i, t := range tags {
		parts[i] = fmt.Sprintf("`%s`", t)
	}
	return strings.Join(parts, ", ")
}
