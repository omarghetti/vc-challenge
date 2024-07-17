package main

import (
	"context"
	"log"

	"github.com/omarghetti/vc-challenge/v2/cmd/server/http"
	"github.com/omarghetti/vc-challenge/v2/internal/api"
	"github.com/omarghetti/vc-challenge/v2/internal/documents"
	"github.com/omarghetti/vc-challenge/v2/internal/repo"
	"github.com/omarghetti/vc-challenge/v2/internal/util"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := util.NewConfig()
	if err != nil {
		log.Fatalf("Error loading configuration file")
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:             config.RedisAddr + ":" + config.RedisPort,
		DB:               0,
		DisableIndentity: true,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Error connecting to Redis")
	}

	//initalize redis repository
	redisRepo := repo.New(rdb)
	// Create a new instance of the documents service
	documents := documents.New(redisRepo)

	// Create a new instance of the API
	api := api.NewServer(documents)

	// Start the HTTP server
	http_service := http.NewService(api, &config)

	defer http_service.Shutdown()

	http_service.Start()

}
