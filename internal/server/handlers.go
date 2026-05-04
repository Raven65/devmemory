package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"devmemory/internal/action"
	"devmemory/internal/core"
	"devmemory/internal/export"
	"devmemory/internal/search"
	"devmemory/internal/store"
)

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleListEntries(w http.ResponseWriter, r *http.Request) {
	opts := store.ListOptions{
		Type:    core.EntryType(r.URL.Query().Get("type")),
		Project: r.URL.Query().Get("project"),
		Tag:     r.URL.Query().Get("tag"),
	}

	entries, err := s.store.List(opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (s *Server) handleCreateEntry(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Type    string   `json:"type"`
		Title   string   `json:"title"`
		Content string   `json:"content"`
		Project string   `json:"project"`
		Tags    []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if input.Content == "" {
		writeError(w, http.StatusBadRequest, "content is required")
		return
	}

	entryType := input.Type
	if entryType == "" {
		entryType = string(action.DetectType(input.Content))
	}

	entry := core.NewEntry(core.EntryType(entryType), input.Content)
	entry.Title = input.Title
	entry.Project = input.Project
	entry.Tags = input.Tags

	if entry.Type == core.EntryTypeCommand && action.IsDangerous(entry.Content) {
		entry.Dangerous = true
	}

	if err := s.store.Create(entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, entry)
}

func (s *Server) handleGetEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	// Try direct ID first, then prefix match
	entry, err := s.store.Get(id)
	if err != nil {
		fullID, resolveErr := s.store.ResolveID(id)
		if resolveErr != nil {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
		entry, err = s.store.Get(fullID)
		if err != nil {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
	}

	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleUpdateEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	entry, err := s.store.Get(id)
	if err != nil {
		fullID, resolveErr := s.store.ResolveID(id)
		if resolveErr != nil {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
		entry, err = s.store.Get(fullID)
		if err != nil {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
	}

	var input struct {
		Type     *string  `json:"type"`
		Title    *string  `json:"title"`
		Content  *string  `json:"content"`
		Project  *string  `json:"project"`
		Tags     []string `json:"tags"`
		Favorite *bool    `json:"favorite"`
		Archived *bool    `json:"archived"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if input.Type != nil {
		entry.Type = core.EntryType(*input.Type)
	}
	if input.Title != nil {
		entry.Title = *input.Title
	}
	if input.Content != nil {
		entry.Content = *input.Content
	}
	if input.Project != nil {
		entry.Project = *input.Project
	}
	if input.Tags != nil {
		entry.Tags = input.Tags
	}
	if input.Favorite != nil {
		entry.Favorite = *input.Favorite
	}
	if input.Archived != nil {
		entry.Archived = *input.Archived
	}

	entry.UpdatedAt = time.Now()

	if err := s.store.Update(entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.store.Delete(id); err != nil {
		// Try prefix match
		fullID, resolveErr := s.store.ResolveID(id)
		if resolveErr != nil {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
		if err := s.store.Delete(fullID); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	entryType := r.URL.Query().Get("type")
	opts := store.ListOptions{
		Type: core.EntryType(entryType),
	}

	entries, err := s.store.List(opts)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	engine := search.NewEngine()
	results := engine.Search(entries, q)

	// Return just the entries from search results, sorted by score
	type searchResponse struct {
		Entry *core.Entry `json:"entry"`
		Score float64     `json:"score"`
	}

	resp := make([]searchResponse, len(results))
	for i, r := range results {
		resp[i] = searchResponse{Entry: r.Entry, Score: r.Score}
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleToday(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	entries, err := s.store.GetByDate(now.Year(), int(now.Month()), now.Day())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"date":    now.Format("2006-01-02"),
		"entries": entries,
	})
}

func (s *Server) handleCopy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	entry, err := s.store.Get(id)
	if err != nil {
		fullID, resolveErr := s.store.ResolveID(id)
		if resolveErr != nil {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
		entry, err = s.store.Get(fullID)
		if err != nil {
			writeError(w, http.StatusNotFound, "entry not found")
			return
		}
	}

	now := time.Now()
	entry.LastUsedAt = &now
	entry.UseCount++
	entry.UpdatedAt = now
	s.store.Update(entry)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"content":   entry.Content,
		"use_count": entry.UseCount,
	})
}

func (s *Server) handleExportToday(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	entries, err := s.store.GetByDate(now.Year(), int(now.Month()), now.Day())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	md, err := export.ExportDailyMarkdown(entries, now)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=devmemory-%s.md", now.Format("2006-01-02")))
	w.Write([]byte(md))
}

func (s *Server) handleExportJSON(w http.ResponseWriter, r *http.Request) {
	entries, err := s.store.List(store.ListOptions{})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data, err := export.ExportJSON(entries)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=devmemory-backup.json")
	w.Write(data)
}

func (s *Server) handleImportJSON(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Entries json.RawMessage `json:"entries"`
	}

	// Try reading as {"entries": [...]} or as direct [...]
	body := json.NewDecoder(r.Body)
	var raw json.RawMessage
	if err := body.Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	entries, err := export.ImportJSON(raw)
	if err != nil {
		// Try wrapped format
		if err2 := json.Unmarshal(raw, &input); err2 == nil {
			entries, err = export.ImportJSON(input.Entries)
		}
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid entry format")
			return
		}
	}

	imported := 0
	for _, entry := range entries {
		if err := s.store.Create(entry); err != nil {
			if err := s.store.Update(entry); err != nil {
				continue
			}
		}
		imported++
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"imported": imported,
		"total":    len(entries),
	})
}

// parseIntOrDefault parses a query parameter as int, returning default on failure.
func parseIntOrDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

// containsString checks if a string contains a substring (case-insensitive).
func containsString(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
