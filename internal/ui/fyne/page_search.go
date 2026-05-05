package fyne

import (

	"devmemory/internal/core"
	"devmemory/internal/search"
	"devmemory/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// searchPage implements the entry search UI.
type searchPage struct {
	svc    *service.MemoryService
	window fyne.Window

	query    *widget.Entry
	typeSel  *widget.Select
	results  *widget.List
	entries  []*search.SearchResult
}

func newSearchPage(svc *service.MemoryService, w fyne.Window) *searchPage {
	p := &searchPage{
		svc:    svc,
		window: w,
	}

	p.query = widget.NewEntry()
	p.query.SetPlaceHolder("Search...")
	p.query.OnSubmitted = func(_ string) { p.doSearch() }

	types := []string{"all"}
	for _, t := range core.ValidEntryTypes() {
		types = append(types, string(t))
	}
	p.typeSel = widget.NewSelect(types, nil)
	p.typeSel.SetSelected("all")

	p.results = widget.NewList(
		func() int { return len(p.entries) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("type"),
				widget.NewLabel("title"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			e := p.entries[i]
			row.Objects[0].(*widget.Label).SetText(string(e.Entry.Type))
			row.Objects[1].(*widget.Label).SetText(truncate(e.Entry.TitleOrContent(), 60))
		},
	)
	p.results.OnSelected = func(i widget.ListItemID) {
		showEntryDialog(p.svc, p.entries[i].Entry, p.window, p.doSearch)
		p.results.Unselect(i)
	}

	return p
}

func (p *searchPage) build() fyne.CanvasObject {
	top := container.NewBorder(
		nil, nil, nil,
		container.NewHBox(p.typeSel, widget.NewButton("Search", p.doSearch)),
		p.query,
	)

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Search Entries", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			top,
		),
		nil, nil, nil,
		p.results,
	)
}

func (p *searchPage) doSearch() {
	q := p.query.Text
	if q == "" {
		return
	}

	var entryType core.EntryType
	if sel := p.typeSel.Selected; sel != "all" {
		entryType = core.EntryType(sel)
	}

	results, err := p.svc.SearchEntries(q, entryType)
	if err != nil {
		return
	}

	p.entries = results
	p.results.Refresh()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
