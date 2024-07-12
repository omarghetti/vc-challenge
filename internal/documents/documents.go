package documents

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Document struct {
	ID           string   `json:"id"`
	Text         string   `json:"text"`
	MatchedWords []string `json:"matched_words"`
}

type Documents struct {
	redis_client *redis.Client
}

func New(redis_client *redis.Client) *Documents {
	return &Documents{
		redis_client: redis_client,
	}
}

func (d *Documents) GetDocByID(ctx context.Context, documentID string) (*Document, error) {
	currCtx, cancel := context.WithTimeout(ctx, 1*time.Second)

	defer cancel()

	val, err := d.redis_client.Get(currCtx, documentID).Result()
	if err != nil {
		return nil, errors.New("document not found")
	}

	return &Document{
		ID:           documentID,
		Text:         val,
		MatchedWords: []string{},
	}, nil

}

func (d *Documents) SetDoc(ctx context.Context, documentID, text string) error {
	currCtx, cancel := context.WithTimeout(ctx, 1*time.Second)

	defer cancel()

	err := d.redis_client.Set(currCtx, documentID, text, 0).Err()
	if err != nil {
		return errors.New("could not set document")
	}

	return nil
}

func (d *Documents) Search(ctx context.Context, query string) ([]Document, error) {
	currCtx, cancel := context.WithTimeout(ctx, 1*time.Second)

	defer cancel()

	keys, err := d.redis_client.Keys(currCtx, "*").Result()
	if err != nil {
		return nil, errors.New("could not search documents")
	}

	var documents []Document
	for _, key := range keys {
		val, err := d.redis_client.Get(currCtx, key).Result()
		if err != nil {
			return nil, errors.New("could not get document")
		}

		documents = append(documents, Document{
			ID:           key,
			Text:         val,
			MatchedWords: []string{},
		})
	}

	return documents, nil
}
