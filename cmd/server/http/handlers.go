package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/omarghetti/vc-challenge/v2/internal/api"
)

type DocumentRequest struct {
	Text string `json:"text"`
}

type GenericErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

func writeJSONErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(GenericErrorResponse{
		StatusCode: code,
		Error:      message,
	})
}

// Handlers is a struct that holds all the handlers for the HTTP server
type Handlers struct {
	apis api.Server
	env  string
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	status, err := h.apis.Health(h.env)
	if err != nil {
		slog.Error(err.Error())
		writeJSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func (h *Handlers) GetDocumentByID(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "documentID")
	document, err := h.apis.GetDocByID(r.Context(), documentID)
	if err != nil {
		slog.Error(err.Error())
		writeJSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *Handlers) SetDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "documentID")

	// max 1MB payload
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	var req DocumentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error(err.Error())
		writeJSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// every document coming in needs to be indexed first and then stored
	err := h.apis.SetDoc(r.Context(), documentID, req.Text)
	if err != nil {
		slog.Error(err.Error())
		writeJSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}

func (h *Handlers) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	documents, err := h.apis.Search(r.Context(), query)
	if err != nil {
		slog.Error(err.Error())
		writeJSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

func (h *Handlers) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	documentID := chi.URLParam(r, "documentID")

	err := h.apis.DeleteDoc(r.Context(), documentID)

	if err != nil {
		slog.Error(err.Error())
		writeJSONErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}
