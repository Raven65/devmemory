package fyne

import (
	"fmt"
	"path/filepath"

	"devmemory/internal/config"
	"devmemory/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// settingsPage shows app info and export/import controls.
type settingsPage struct {
	svc    *service.MemoryService
	window fyne.Window
	info   *widget.Label
}

func newSettingsPage(svc *service.MemoryService, w fyne.Window) *settingsPage {
	return &settingsPage{
		svc:    svc,
		window: w,
	}
}

func (p *settingsPage) build() fyne.CanvasObject {
	cfg := config.DefaultConfig()
	count, _ := p.svc.Count()

	p.info = widget.NewLabelWithStyle(
		fmt.Sprintf("Version:     dev\nData Dir:    %s\nDatabase:    %s\nTotal Entries: %d",
			cfg.DataDir, filepath.Base(cfg.DBPath), count),
		fyne.TextAlignLeading, fyne.TextStyle{Monospace: true},
	)

	exportTodayBtn := widget.NewButton("Export Today (Markdown)", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			md, err := p.svc.ExportTodayMarkdown()
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			writer.Write([]byte(md))
		}, p.window)
		fd.SetFileName("devmemory-today.md")
		fd.Show()
	})

	exportJSONBtn := widget.NewButton("Export JSON", func() {
		fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer writer.Close()
			data, err := p.svc.ExportJSON()
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			writer.Write(data)
		}, p.window)
		fd.SetFileName("devmemory-export.json")
		fd.Show()
	})

	importJSONBtn := widget.NewButton("Import JSON", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()
			data := make([]byte, 1024*1024) // 1MB limit
			n, err := reader.Read(data)
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			imported, err := p.svc.ImportJSON(data[:n])
			if err != nil {
				dialog.ShowError(err, p.window)
				return
			}
			dialog.ShowInformation("Import", fmt.Sprintf("Imported %d entries.", imported), p.window)
			p.refreshInfo()
		}, p.window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		fd.Show()
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Settings", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		p.info,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("Data Management", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		exportTodayBtn,
		exportJSONBtn,
		importJSONBtn,
	)
}

func (p *settingsPage) refreshInfo() {
	count, _ := p.svc.Count()
	cfg := config.DefaultConfig()
	p.info.SetText(
		fmt.Sprintf("Version:     dev\nData Dir:    %s\nDatabase:    %s\nTotal Entries: %d",
			cfg.DataDir, filepath.Base(cfg.DBPath), count),
	)
}
