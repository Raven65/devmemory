package fyne

import (
	"devmemory/internal/core"
	"devmemory/internal/service"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// capturePage implements the entry creation form.
type capturePage struct {
	svc    *service.MemoryService
	window fyne.Window

	content  *widget.Entry
	title    *widget.Entry
	project  *widget.Entry
	tags     *widget.Entry
	typeSel  *widget.Select
}

func newCapturePage(svc *service.MemoryService, w fyne.Window) *capturePage {
	p := &capturePage{
		svc:    svc,
		window: w,
	}

	p.content = widget.NewMultiLineEntry()
	p.content.SetPlaceHolder("Enter content here...")

	p.title = widget.NewEntry()
	p.title.SetPlaceHolder("(optional)")

	p.project = widget.NewEntry()
	p.project.SetPlaceHolder("(optional)")

	p.tags = widget.NewEntry()
	p.tags.SetPlaceHolder("comma-separated (optional)")

	// Auto-detect + all entry types
	types := []string{"auto-detect"}
	for _, t := range core.ValidEntryTypes() {
		types = append(types, string(t))
	}
	p.typeSel = widget.NewSelect(types, nil)
	p.typeSel.SetSelected("auto-detect")

	return p
}

func (p *capturePage) build() fyne.CanvasObject {
	form := container.NewVBox(
		widget.NewLabelWithStyle("New Entry", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			&widget.FormItem{Text: "Content", Widget: p.content},
			&widget.FormItem{Text: "Type", Widget: p.typeSel},
			&widget.FormItem{Text: "Title", Widget: p.title},
			&widget.FormItem{Text: "Project", Widget: p.project},
			&widget.FormItem{Text: "Tags", Widget: p.tags},
		),
		layout.NewSpacer(),
		widget.NewButton("Submit", p.submit),
	)
	return container.NewPadded(form)
}

func (p *capturePage) submit() {
	content := strings.TrimSpace(p.content.Text)
	if content == "" {
		dialog.ShowInformation("Validation", "Content is required.", p.window)
		return
	}

	var entryType core.EntryType
	if sel := p.typeSel.Selected; sel != "auto-detect" {
		entryType = core.EntryType(sel)
	}

	var tags []string
	if t := strings.TrimSpace(p.tags.Text); t != "" {
		for _, tag := range strings.Split(t, ",") {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tags = append(tags, tag)
			}
		}
	}

	entry, err := p.svc.CreateEntry(service.CreateEntryInput{
		Content: content,
		Type:    entryType,
		Title:   strings.TrimSpace(p.title.Text),
		Project: strings.TrimSpace(p.project.Text),
		Tags:    tags,
	})
	if err != nil {
		dialog.ShowError(err, p.window)
		return
	}

	dialog.ShowInformation("Created",
		"Entry created: "+entry.ID[:8]+" ("+string(entry.Type)+")",
		p.window,
	)

	// Clear form
	p.content.SetText("")
	p.title.SetText("")
	p.project.SetText("")
	p.tags.SetText("")
	p.typeSel.SetSelected("auto-detect")
}
