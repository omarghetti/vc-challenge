package documents

import (
	"context"
	"errors"
	"strings"

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
	val, err := d.redis_client.Get(ctx, documentID).Result()
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
	// we take into consideration that no parsing is necessary to
	// extract the words from the text
	textList := strings.Split(text, " ")
	_, fail := d.redis_client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, word := range textList {
			// we are only interested in words that are longer than 3 characters
			// we actually omit articles and propositions, not worth indexing
			if len(word) > 3 {
				// we are assuming that the word is valid
				word = strings.ToLower(word)
				pipe.SAdd(ctx, word, documentID)
			}
		}
		return nil
	})

	if fail != nil {
		return errors.New("could not set document, error while indexing words")
	}

	err := d.redis_client.Set(ctx, documentID, text, 0).Err()
	if err != nil {
		return errors.New("could not set document, error while saving document")
	}

	return nil
}

func (d *Documents) Search(ctx context.Context, query string) ([]Document, error) {
	query = strings.ToLower(query)
	words := strings.Split(query, ",")

	docsIds := make(map[string]int)
	var documents []Document

	// we are going to query redis to get set membership for all the words
	// in the query, then we are going to intersect the results and return
	// the documents that contain all the words in the query
	cmds, err := d.redis_client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, word := range words {
			word = strings.ToLower(word)
			pipe.SMembers(ctx, word)
		}
		return nil
	})

	if err != nil {
		return nil, errors.New("could not search documents, error while querying words")
	}

	for _, cmd := range cmds {
		for _, docID := range cmd.(*redis.StringSliceCmd).Val() {
			docsIds[docID]++
		}
	}

	for docID := range docsIds {
		doc, err := d.GetDocByID(ctx, docID)
		if err != nil {
			return nil, errors.New("could not search documents, error while getting document")
		}
		documents = append(documents, *doc)
	}

	return documents, nil
}
