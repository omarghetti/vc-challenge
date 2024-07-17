package documents

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/omarghetti/vc-challenge/v2/internal/repo"
)

type SearchResult struct {
	Query string     `json:"query"`
	Docs  []Document `json:"docs"`
}

type Document struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Documents struct {
	repository repo.Storer
}

func New(storer repo.Storer) *Documents {
	return &Documents{
		repository: storer,
	}
}

func (d *Documents) GetDocByID(ctx context.Context, documentID string) (*Document, error) {
	val, err := d.repository.GetDoc(ctx, documentID)
	if err != nil {
		return nil, errors.New("document not found")
	}

	return &Document{
		ID:   documentID,
		Text: val,
	}, nil

}

func (d *Documents) SetDoc(ctx context.Context, documentID, text string) error {
	// we take into consideration that no parsing is necessary to
	// extract the words from the text
	err := d.repository.SetNewDoc(ctx, documentID, text)
	if err != nil {
		return err
	}

	return nil
}

func (d *Documents) Search(ctx context.Context, query string) (SearchResult, error) {
	var wg sync.WaitGroup
	var result []Document
	query = strings.ToLower(query)
	documents, err := d.repository.SearchDocs(ctx, query)

	if err != nil {
		return SearchResult{}, err
	}

	for _, doc := range documents {
		wg.Add(1)
		go func(doc string) {
			defer wg.Done()
			text, err := d.GetDocByID(ctx, doc)
			if err != nil {
				return
			}
			result = append(result, *text)
		}(doc)
		wg.Wait()
	}

	sr := SearchResult{
		Query: query,
		Docs:  result,
	}

	return sr, nil
}

func (d *Documents) DeleteDoc(ctx context.Context, documentID string) error {
	err := d.repository.DelDoc(ctx, documentID)
	if err != nil {
		return err
	}

	return nil
}
