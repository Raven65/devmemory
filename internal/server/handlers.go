package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"devmemory/internal/core"
	"devmemory/internal/service"
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

	entries, err := s.svc.ListEntries(opts)
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

	entry, err := s.svc.CreateEntry(service.CreateEntryInput{
		Content: input.Content,
		Type:    core.EntryType(input.Type),
		Title:   input.Title,
		Project: input.Project,
		Tags:    input.Tags,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, entry)
}

func (s *Server) handleGetEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	entry, err := s.svc.GetEntry(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "entry not found")
		return
	}

	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleUpdateEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

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

	fields := service.UpdateEntryFields{}
	if input.Type != nil {
		t := core.EntryType(*input.Type)
		fields.Type = &t
	}
	if input.Title != nil {
		fields.Title = input.Title
	}
	if input.Content != nil {
		fields.Content = input.Content
	}
	if input.Project != nil {
		fields.Project = input.Project
	}
	if input.Tags != nil {
		fields.Tags = &input.Tags
	}
	if input.Favorite != nil {
		fields.Favorite = input.Favorite
	}
	if input.Archived != nil {
		fields.Archived = input.Archived
	}

	entry, err := s.svc.UpdateEntry(id, fields)
	if err != nil {
		writeError(w, http.StatusNotFound, "entry not found or update failed")
		return
	}

	writeJSON(w, http.StatusOK, entry)
}

func (s *Server) handleDeleteEntry(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := s.svc.DeleteEntry(id); err != nil {
		writeError(w, http.StatusNotFound, "entry not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		writeError(w, http.StatusBadRequest, "query parameter 'q' is required")
		return
	}

	entryType := core.EntryType(r.URL.Query().Get("type"))

	results, err := s.svc.SearchEntries(q, entryType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

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
	entries, err := s.svc.GetToday()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"date":    time.Now().Format("2006-01-02"),
		"entries": entries,
	})
}

func (s *Server) handleCopy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	entry, err := s.svc.CopyEntry(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "entry not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"content":   entry.Content,
		"use_count": entry.UseCount,
	})
}

func (s *Server) handleExportToday(w http.ResponseWriter, r *http.Request) {
	md, err := s.svc.ExportTodayMarkdown()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dateStr := time.Now().Format("2006-01-02")

	w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=devmemory-%s.md", dateStr))
	w.Write([]byte(md))
}

func (s *Server) handleExportJSON(w http.ResponseWriter, r *http.Request) {
	data, err := s.svc.ExportJSON()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=devmemory-backup.json")
	w.Write(data)
}

func (s *Server) handleImportJSON(w http.ResponseWriter, r *http.Request) {
	var raw json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	imported, err := s.svc.ImportJSON(raw)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid entry format")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"imported": imported,
	})
}
