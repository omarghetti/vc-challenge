package repo

import (
	"context"
	"errors"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Storer interface {
	SetNewDoc(ctx context.Context, key, value string) error
	GetDoc(ctx context.Context, key string) (string, error)
	DelDoc(ctx context.Context, key string) error
	SearchDocs(ctx context.Context, query string) ([]string, error)
}

type RedisStorer struct {
	redis_client *redis.Client
}

func New(redis_client *redis.Client) *RedisStorer {
	return &RedisStorer{
		redis_client: redis_client,
	}
}

// SetNewDoc stores a new document in the redis database, and then
// indexes the words in the document. if document is already in the database
// returns an error
func (r *RedisStorer) SetNewDoc(ctx context.Context, key, value string) error {
	_, ok := r.redis_client.Get(ctx, key).Result()
	if ok == nil {
		return errors.New("document already exists")
	}

	words := strings.Split(value, " ")
	_, fail := r.redis_client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, word := range words {
			if len(word) > 3 {
				word = strings.ToLower(word)
				pipe.SAdd(ctx, word, key)
			}
		}
		return nil
	})

	if fail != nil {
		return errors.New(fail.Error())
	}

	_, err := r.redis_client.Set(ctx, key, value, 0).Result()

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// GetDoc retrieves a document from the redis database given a key
// if the document is not found, an error is returned
func (r *RedisStorer) GetDoc(ctx context.Context, key string) (string, error) {
	return r.redis_client.Get(ctx, key).Result()
}

// DelDoc deletes a document from the redis database given a key and sanitizes the indexes
// if the document is not found, an error is returned
func (r *RedisStorer) DelDoc(ctx context.Context, key string) error {
	val, ok := r.redis_client.Get(ctx, key).Result()
	if ok != nil {
		return errors.New("document not found to delete")
	}

	words := strings.Split(val, " ")
	_, fail := r.redis_client.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, word := range words {
			if len(word) > 3 {
				word = strings.ToLower(word)
				pipe.SRem(ctx, word, key)
			}
		}
		return nil
	})

	if fail != nil {
		return errors.New(fail.Error())
	}

	_, err := r.redis_client.Del(ctx, key).Result()

	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// SearchDocs searches for documents in the redis database that contains all the words
// in the query. The search is case insensitive and the words are separated by commas.
func (r *RedisStorer) SearchDocs(ctx context.Context, query string) ([]string, error) {
	words := strings.Split(query, ",")

	list, err := r.redis_client.SInter(ctx, words...).Result()

	if err != nil {
		return nil, errors.New(err.Error())
	}

	return list, nil

}
