package http

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/omarghetti/vc-challenge/v2/internal/api"
	"github.com/omarghetti/vc-challenge/v2/internal/util"
)

type HTTP struct {
	server *http.Server
}

func (h *HTTP) Start() {
	h.server.ListenAndServe()
}

func (h *HTTP) Shutdown() {
	h.server.Shutdown(context.TODO())
}

func NewService(api api.Server, config *util.Config) *HTTP {
	r := chi.NewRouter()
	h := &Handlers{
		apis: api,
		env:  config.Environment,
	}

	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/health", h.Health)
		r.Get("/document/{documentID}", h.GetDocumentByID)
		r.Post("/document/{documentID}", h.SetDocument)
		r.Get("/search", h.Search)
		r.Delete("/document/{documentID}", h.DeleteDocument)
	})

	server := &http.Server{
		Addr:    config.HTTPAddr,
		Handler: r,
	}

	return &HTTP{
		server: server,
	}
}
