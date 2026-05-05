package fyne

import (
	"fmt"

	"devmemory/internal/config"
	"devmemory/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// pageDef describes a sidebar navigation item.
type pageDef struct {
	name   string
	icon   fyne.Resource
	build  func() fyne.CanvasObject
}

// mainWindow is the primary application window.
type mainWindow struct {
	fyne.Window
	app       fyne.App
	svc       *service.MemoryService
	pages     []pageDef
	statusBar *widget.Label
}

func newMainWindow(a fyne.App, svc *service.MemoryService) *mainWindow {
	w := a.NewWindow("DevMemory")
	w.Resize(fyne.NewSize(900, 600))

	mw := &mainWindow{
		Window: w,
		app:    a,
		svc:    svc,
	}

	mw.pages = []pageDef{
		{name: "Capture", icon: theme.ContentAddIcon(), build: func() fyne.CanvasObject {
			return newCapturePage(svc, w).build()
		}},
		{name: "Search", icon: theme.SearchIcon(), build: func() fyne.CanvasObject {
			return newSearchPage(svc, w).build()
		}},
		{name: "Today", icon: theme.CalendarIcon(), build: func() fyne.CanvasObject {
			return newTodayPage(svc, w).build()
		}},
		{name: "Actions", icon: theme.MediaPlayIcon(), build: func() fyne.CanvasObject {
			return newActionsPage(svc, w).build()
		}},
		{name: "Knowledge", icon: theme.DocumentIcon(), build: func() fyne.CanvasObject {
			return newKnowledgePage(svc, w).build()
		}},
		{name: "Settings", icon: theme.SettingsIcon(), build: func() fyne.CanvasObject {
			return newSettingsPage(svc, w).build()
		}},
	}

	// Build sidebar
	navList := widget.NewList(
		func() int { return len(mw.pages) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.ContentAddIcon()),
				widget.NewLabel("Knowledge   "),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			row := o.(*fyne.Container)
			row.Objects[0].(*widget.Icon).SetResource(mw.pages[i].icon)
			row.Objects[1].(*widget.Label).SetText(mw.pages[i].name)
		},
	)

	// Eagerly build all page content
	pageObjects := make([]fyne.CanvasObject, len(mw.pages))
	for i, p := range mw.pages {
		pageObjects[i] = p.build()
	}

	contentArea := container.NewStack(pageObjects...)
	navList.OnSelected = func(i widget.ListItemID) {
		for idx, obj := range pageObjects {
			if idx == i {
				obj.Show()
			} else {
				obj.Hide()
			}
		}
		contentArea.Refresh()
		mw.refreshStatus()
	}
	navList.Select(0)

	// Status bar
	mw.statusBar = widget.NewLabel(mw.statusText())
	mw.statusBar.TextStyle = fyne.TextStyle{Italic: true}

	// Sidebar with app title
	sidebar := container.NewBorder(
		widget.NewLabelWithStyle("DevMemory", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		navList,
	)

	// Main layout: sidebar | content, status bar at bottom
	body := container.NewBorder(nil, nil, sidebar, nil, contentArea)

	w.SetContent(container.NewBorder(
		nil,
		container.NewHBox(layout.NewSpacer(), mw.statusBar),
		nil, nil,
		body,
	))

	return mw
}

func (mw *mainWindow) statusText() string {
	count, err := mw.svc.Count()
	if err != nil {
		count = 0
	}
	cfg := config.DefaultConfig()
	return fmt.Sprintf("v0.1 | %s | %d entries", cfg.DataDir, count)
}

func (mw *mainWindow) refreshStatus() {
	if mw.statusBar != nil {
		mw.statusBar.SetText(mw.statusText())
	}
}
