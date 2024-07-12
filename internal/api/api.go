package api

import (
	"context"

	"github.com/omarghetti/vc-challenge/v2/internal/documents"
)

type Server interface {
	GetDocByID(ctx context.Context, documentID string) (*documents.Document, error)
	SetDoc(ctx context.Context, documentID, text string) error
	Search(ctx context.Context, query string) ([]documents.Document, error)
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
		"env":    env,
		"status": "ok",
	}, nil
}

func (a *API) GetDocByID(ctx context.Context, documentID string) (*documents.Document, error) {
	return a.documents.GetDocByID(ctx, documentID)
}

func (a *API) SetDoc(ctx context.Context, documentID, text string) error {
	return a.documents.SetDoc(ctx, documentID, text)
}

func (a *API) Search(ctx context.Context, query string) ([]documents.Document, error) {
	return a.documents.Search(ctx, query)
}
