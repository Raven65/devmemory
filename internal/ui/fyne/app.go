package fyne

import (
	"devmemory/internal/config"
	"devmemory/internal/service"
	"devmemory/internal/store"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// App wraps the Fyne application and its dependencies.
type App struct {
	fyneApp fyne.App
	svc     *service.MemoryService
	cleanup func()
}

// NewApp creates a new Fyne application with an initialized service layer.
func NewApp() *App {
	cfg := config.DefaultConfig()
	if err := config.EnsureDataDir(cfg.DataDir); err != nil {
		panic(err)
	}

	s := store.NewBBoltStore(cfg.DBPath)
	if err := s.Open(); err != nil {
		panic(err)
	}

	return &App{
		fyneApp: app.NewWithID("devmemory"),
		svc:     service.NewMemoryService(s),
		cleanup: func() { s.Close() },
	}
}

// Run starts the Fyne event loop.
func (a *App) Run() {
	w := newMainWindow(a.fyneApp, a.svc)
	w.ShowAndRun()
	a.cleanup()
}

// Close shuts down the application.
func (a *App) Close() {
	a.cleanup()
}
