package api

import (
	"context"
	"time"

	"github.com/omarghetti/vc-challenge/v2/internal/documents"
)

var (
	now = time.Now()
)

type Server interface {
	GetDocByID(ctx context.Context, documentID string) (*documents.Document, error)
	SetDoc(ctx context.Context, documentID, text string) error
	Search(ctx context.Context, query string) (documents.SearchResult, error)
	DeleteDoc(ctx context.Context, documentID string) error
	Health(env string) (map[string]any, error)
}

type API struct {
	documents *documents.Documents
}

func New(documents *documents.Documents) *API {
	return &API{
		documents: documents,
	}
}

func NewServer(documents *documents.Documents) Server {
	return New(documents)
}

func (a *API) Health(env string) (map[string]any, error) {
	return map[string]any{
		"env":          env,
		"status":       "ok",
		"current_time": now.String(),
	}, nil
}

func (a *API) GetDocByID(ctx context.Context, documentID string) (*documents.Document, error) {
	newCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	return a.documents.GetDocByID(newCtx, documentID)
}

func (a *API) SetDoc(ctx context.Context, documentID, text string) error {
	newCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	return a.documents.SetDoc(newCtx, documentID, text)
}

func (a *API) Search(ctx context.Context, query string) (documents.SearchResult, error) {
	newCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	return a.documents.Search(newCtx, query)
}

func (a *API) DeleteDoc(ctx context.Context, documentID string) error {
	newCtx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()
	return a.documents.DeleteDoc(newCtx, documentID)
}
