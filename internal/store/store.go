package store

import "devmemory/internal/core"

// ListOptions controls filtering when listing entries.
type ListOptions struct {
	Type     core.EntryType
	Project  string
	Tag      string
	Archived bool
	Limit    int
	Offset   int
}

// Store is the interface for entry persistence.
type Store interface {
	Create(entry *core.Entry) error
	Get(id string) (*core.Entry, error)
	Update(entry *core.Entry) error
	Delete(id string) error
	List(opts ListOptions) ([]*core.Entry, error)
	GetByDate(year, month, day int) ([]*core.Entry, error)
	ResolveID(prefix string) (string, error)
	Open() error
	Close() error
}
