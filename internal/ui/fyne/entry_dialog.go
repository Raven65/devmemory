package fyne

import (
	"fmt"
	"strings"
	"time"

	"devmemory/internal/core"
	"devmemory/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// showEntryDialog displays a dialog with entry details and action buttons.
func showEntryDialog(svc *service.MemoryService, entry *core.Entry, parent fyne.Window, onModified func()) {
	title := entry.TitleOrContent()
	if len(title) > 40 {
		title = title[:40] + "..."
	}

	// Detail labels
	info := widget.NewLabel(formatEntryDetail(entry))
	info.Wrapping = fyne.TextWrapWord

	// Action buttons
	copyBtn := widget.NewButton("Copy", func() {
		if _, err := svc.CopyEntry(entry.ID); err != nil {
			dialog.ShowError(err, parent)
			return
		}
		dialog.ShowInformation("Copied", "Content copied. Use count updated.", parent)
		if onModified != nil {
			onModified()
		}
	})

	favBtn := widget.NewButton("Toggle Favorite", func() {
		if _, err := svc.ToggleFavorite(entry.ID); err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if onModified != nil {
			onModified()
		}
	})

	archiveBtn := widget.NewButton("Toggle Archive", func() {
		if _, err := svc.ToggleArchive(entry.ID); err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if onModified != nil {
			onModified()
		}
	})

	editBtn := widget.NewButton("Edit", func() {
		showEditDialog(svc, entry, parent, onModified)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		dialog.ShowConfirm("Delete Entry",
			fmt.Sprintf("Delete [%s] %s?", entry.Type, entry.TitleOrContent()),
			func(ok bool) {
				if !ok {
					return
				}
				if err := svc.DeleteEntry(entry.ID); err != nil {
					dialog.ShowError(err, parent)
					return
				}
				if onModified != nil {
					onModified()
				}
			}, parent)
	})
	deleteBtn.Importance = widget.HighImportance

	actions := container.NewHBox(copyBtn, favBtn, archiveBtn, editBtn, deleteBtn)

	content := container.NewVBox(
		info,
		widget.NewSeparator(),
		actions,
	)

	d := dialog.NewCustom(title, "Close", content, parent)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}

// showEditDialog shows a dialog for editing entry fields.
func showEditDialog(svc *service.MemoryService, entry *core.Entry, parent fyne.Window, onModified func()) {
	titleEntry := widget.NewEntry()
	titleEntry.SetText(entry.Title)

	contentEntry := widget.NewMultiLineEntry()
	contentEntry.SetText(entry.Content)

	types := make([]string, len(core.ValidEntryTypes()))
	for i, t := range core.ValidEntryTypes() {
		types[i] = string(t)
	}
	typeSel := widget.NewSelect(types, nil)
	typeSel.SetSelected(string(entry.Type))

	projectEntry := widget.NewEntry()
	projectEntry.SetText(entry.Project)

	tagsEntry := widget.NewEntry()
	tagsEntry.SetText(strings.Join(entry.Tags, ", "))

	form := container.NewVBox(
		widget.NewForm(
			&widget.FormItem{Text: "Title", Widget: titleEntry},
			&widget.FormItem{Text: "Content", Widget: contentEntry},
			&widget.FormItem{Text: "Type", Widget: typeSel},
			&widget.FormItem{Text: "Project", Widget: projectEntry},
			&widget.FormItem{Text: "Tags", Widget: tagsEntry},
		),
	)

	dialog.ShowCustomConfirm("Edit Entry", "Save", "Cancel", form, func(save bool) {
		if !save {
			return
		}
		t := core.EntryType(typeSel.Selected)
		tags := splitTagsStr(tagsEntry.Text)
		_, err := svc.UpdateEntry(entry.ID, service.UpdateEntryFields{
			Title:   &titleEntry.Text,
			Content: &contentEntry.Text,
			Type:    &t,
			Project: &projectEntry.Text,
			Tags:    &tags,
		})
		if err != nil {
			dialog.ShowError(err, parent)
			return
		}
		if onModified != nil {
			onModified()
		}
	}, parent)
}

func formatEntryDetail(e *core.Entry) string {
	var b strings.Builder
	fmt.Fprintf(&b, "ID:        %s\n", e.ID)
	fmt.Fprintf(&b, "Type:      %s\n", e.Type)
	if e.Title != "" {
		fmt.Fprintf(&b, "Title:     %s\n", e.Title)
	}
	fmt.Fprintf(&b, "Content:   %s\n", e.Content)
	if e.Summary != "" {
		fmt.Fprintf(&b, "Summary:   %s\n", e.Summary)
	}
	if e.Project != "" {
		fmt.Fprintf(&b, "Project:   %s\n", e.Project)
	}
	if len(e.Tags) > 0 {
		fmt.Fprintf(&b, "Tags:      %s\n", strings.Join(e.Tags, ", "))
	}
	fmt.Fprintf(&b, "Favorite:  %v\n", e.Favorite)
	fmt.Fprintf(&b, "Dangerous: %v\n", e.Dangerous)
	fmt.Fprintf(&b, "Archived:  %v\n", e.Archived)
	fmt.Fprintf(&b, "UseCount:  %d\n", e.UseCount)
	fmt.Fprintf(&b, "Created:   %s\n", e.CreatedAt.Format(time.DateTime))
	fmt.Fprintf(&b, "Updated:   %s\n", e.UpdatedAt.Format(time.DateTime))
	if e.LastUsedAt != nil {
		fmt.Fprintf(&b, "LastUsed:  %s\n", e.LastUsedAt.Format(time.DateTime))
	}
	return b.String()
}

func splitTagsStr(s string) []string {
	parts := strings.Split(s, ",")
	tags := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			tags = append(tags, p)
		}
	}
	return tags
}
