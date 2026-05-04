package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"devmemory/internal/core"

	bbolt "go.etcd.io/bbolt"
)

const bucketName = "entries"

// BBoltStore implements Store using bbolt.
type BBoltStore struct {
	dbPath string
	db     *bbolt.DB
}

// NewBBoltStore creates a new BBoltStore pointing to the given database file.
func NewBBoltStore(dbPath string) *BBoltStore {
	return &BBoltStore{dbPath: dbPath}
}

// Open creates the data directory and opens the database.
func (s *BBoltStore) Open() error {
	dir := filepath.Dir(s.dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}

	db, err := bbolt.Open(s.dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	s.db = db

	return s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

// Close closes the database.
func (s *BBoltStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Create stores a new entry.
func (s *BBoltStore) Create(entry *core.Entry) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("marshal entry: %w", err)
		}
		return b.Put([]byte(entry.ID), data)
	})
}

// Get retrieves an entry by ID.
func (s *BBoltStore) Get(id string) (*core.Entry, error) {
	var entry *core.Entry
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		data := b.Get([]byte(id))
		if data == nil {
			return fmt.Errorf("entry not found: %s", id)
		}
		var e core.Entry
		if err := json.Unmarshal(data, &e); err != nil {
			return fmt.Errorf("unmarshal entry: %w", err)
		}
		entry = &e
		return nil
	})
	return entry, err
}

// Update overwrites an existing entry.
func (s *BBoltStore) Update(entry *core.Entry) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b.Get([]byte(entry.ID)) == nil {
			return fmt.Errorf("entry not found: %s", entry.ID)
		}
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("marshal entry: %w", err)
		}
		return b.Put([]byte(entry.ID), data)
	})
}

// Delete removes an entry by ID.
func (s *BBoltStore) Delete(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b.Get([]byte(id)) == nil {
			return fmt.Errorf("entry not found: %s", id)
		}
		return b.Delete([]byte(id))
	})
}

// List returns entries matching the given filter options, sorted by CreatedAt desc.
func (s *BBoltStore) List(opts ListOptions) ([]*core.Entry, error) {
	var entries []*core.Entry
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e core.Entry
			if err := json.Unmarshal(v, &e); err != nil {
				continue // skip malformed entries
			}

			if !matchFilter(&e, opts) {
				continue
			}

			entries = append(entries, &e)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt.After(entries[j].CreatedAt)
	})

	return applyPagination(entries, opts), nil
}

// GetByDate returns entries created on the given date (year, month, day).
func (s *BBoltStore) GetByDate(year, month, day int) ([]*core.Entry, error) {
	start := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	end := start.Add(24 * time.Hour)

	var entries []*core.Entry
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e core.Entry
			if err := json.Unmarshal(v, &e); err != nil {
				continue
			}
			if !e.CreatedAt.Before(start) && e.CreatedAt.Before(end) {
				entries = append(entries, &e)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].CreatedAt.After(entries[j].CreatedAt)
	})

	return entries, nil
}

// ResolveID resolves a short ID prefix to a full ID.
// Returns error if no match or multiple matches found.
func (s *BBoltStore) ResolveID(prefix string) (string, error) {
	var matches []string
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			id := string(k)
			if strings.HasPrefix(id, prefix) {
				matches = append(matches, id)
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no entry matching prefix: %s", prefix)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("ambiguous prefix %s matches %d entries", prefix, len(matches))
	}
}

func matchFilter(e *core.Entry, opts ListOptions) bool {
	if opts.Type != "" && e.Type != opts.Type {
		return false
	}
	if opts.Project != "" && !strings.EqualFold(e.Project, opts.Project) {
		return false
	}
	if opts.Tag != "" && !containsTag(e.Tags, opts.Tag) {
		return false
	}
	if !opts.Archived && e.Archived {
		return false
	}
	return true
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if strings.EqualFold(t, tag) {
			return true
		}
	}
	return false
}

func applyPagination(entries []*core.Entry, opts ListOptions) []*core.Entry {
	if opts.Offset > 0 {
		if opts.Offset >= len(entries) {
			return nil
		}
		entries = entries[opts.Offset:]
	}
	if opts.Limit > 0 && opts.Limit < len(entries) {
		entries = entries[:opts.Limit]
	}
	return entries
}
