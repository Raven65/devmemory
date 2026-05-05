package fyne

import (
	"devmemory/internal/core"
	"devmemory/internal/service"
	"devmemory/internal/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// actionTypes are the entry types shown on the Actions page.
var actionTypes = []core.EntryType{
	core.EntryTypeCommand, core.EntryTypeURL, core.EntryTypeSnippet,
	core.EntryTypePrompt, core.EntryTypeFile, core.EntryTypeFolder,
}

// actionsPage shows action-type entries with Copy/Open/Execute buttons.
type actionsPage struct {
	svc     *service.MemoryService
	window  fyne.Window
	entries []*EntryBrief
	list    *widget.List
}

func newActionsPage(svc *service.MemoryService, w fyne.Window) *actionsPage {
	p := &actionsPage{
		svc:    svc,
		window: w,
	}

	p.list = widget.NewList(
		func() int { return len(p.entries) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("[type]"),
				widget.NewLabel("Title"),
				widget.NewLabel("(project)"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			e := p.entries[i]
			row.Objects[0].(*widget.Label).SetText("[" + e.Type + "]")
			row.Objects[1].(*widget.Label).SetText(truncate(e.Title, 40))
			row.Objects[2].(*widget.Label).SetText(e.Project)
		},
	)
	p.list.OnSelected = func(i widget.ListItemID) {
		entry, err := p.svc.GetEntry(p.entries[i].ID)
		if err != nil {
			return
		}
		showEntryDialog(p.svc, entry, p.window, p.refresh)
		p.list.Unselect(i)
	}

	return p
}

func (p *actionsPage) build() fyne.CanvasObject {
	header := container.NewHBox(
		widget.NewLabelWithStyle("Actions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewButton("Refresh", p.refresh),
	)

	p.refresh()

	return container.NewBorder(header, nil, nil, nil, p.list)
}

func (p *actionsPage) refresh() {
	var allEntries []*EntryBrief
	for _, t := range actionTypes {
		entries, err := p.svc.ListEntries(store.ListOptions{Type: t})
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.Archived {
				continue
			}
			allEntries = append(allEntries, &EntryBrief{
				ID:      e.ID,
				Type:    string(e.Type),
				Title:   e.TitleOrContent(),
				Project: e.Project,
				Tags:    e.Tags,
			})
		}
	}
	p.entries = allEntries
	p.list.Refresh()
}
