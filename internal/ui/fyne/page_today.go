package fyne

import (
	"fmt"
	"strings"
	"time"

	"devmemory/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// todayPage shows today's entries as a timeline.
type todayPage struct {
	svc    *service.MemoryService
	window fyne.Window
	list   *widget.List
	entries []*EntryBrief
}

// EntryBrief holds display data for a list entry.
type EntryBrief struct {
	ID      string
	Type    string
	Title   string
	Project string
	Tags    []string
	Time    string
}

func newTodayPage(svc *service.MemoryService, w fyne.Window) *todayPage {
	p := &todayPage{
		svc:    svc,
		window: w,
	}

	p.list = widget.NewList(
		func() int { return len(p.entries) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("00:00"),
				widget.NewLabel("[type]"),
				widget.NewLabel("Title"),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			e := p.entries[i]
			row.Objects[0].(*widget.Label).SetText(e.Time)
			row.Objects[1].(*widget.Label).SetText("[" + e.Type + "]")
			row.Objects[2].(*widget.Label).SetText(e.Title)
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

func (p *todayPage) build() fyne.CanvasObject {
	header := container.NewHBox(
		widget.NewLabelWithStyle(
			fmt.Sprintf("Today (%s)", time.Now().Format("2006-01-02")),
			fyne.TextAlignLeading, fyne.TextStyle{Bold: true},
		),
		widget.NewButton("Refresh", p.refresh),
	)

	p.refresh()

	return container.NewBorder(header, nil, nil, nil, p.list)
}

func (p *todayPage) refresh() {
	entries, err := p.svc.GetToday()
	if err != nil {
		p.entries = nil
		p.list.Refresh()
		return
	}

	p.entries = make([]*EntryBrief, len(entries))
	for i, e := range entries {
		brief := &EntryBrief{
			ID:      e.ID,
			Type:    string(e.Type),
			Title:   e.TitleOrContent(),
			Project: e.Project,
			Tags:    e.Tags,
			Time:    e.CreatedAt.Format("15:04"),
		}
		p.entries[i] = brief
	}
	p.list.Refresh()
}

// entryBriefSubtitle builds an optional subtitle string.
func entryBriefSubtitle(e *EntryBrief) string {
	var parts []string
	if e.Project != "" {
		parts = append(parts, e.Project)
	}
	if len(e.Tags) > 0 {
		parts = append(parts, "#"+strings.Join(e.Tags, " #"))
	}
	return strings.Join(parts, " ")
}
