package main

import (
	"context"
	"go-workshop-example/config"
	"go-workshop-example/pkg/adapter"
	"go-workshop-example/pkg/repository"
	"go-workshop-example/pkg/usecase"
	"log"
	"time"
)

func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongodbClient, err := adapter.NewMongoDBConnection(ctx, cfg.MongoDBURI)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = mongodbClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	mongoDBAdapter := adapter.NewMongoDBAdapter(mongodbClient)
	hashtagCollection := mongodbClient.Database("go-workshop").Collection("hashtags")

	rabbitmqAdapter := adapter.NewRabbitMQAdapter()
	defer rabbitmqAdapter.CleanUp()

	hashtagRepo := repository.NewHashtag(mongoDBAdapter, hashtagCollection, rabbitmqAdapter)

	hashtagUsecase := usecase.NewHashtag(hashtagRepo)

	hashtags, err := hashtagUsecase.GetHashtags(ctx)
	if err != nil {
		panic(err)
	}

	err = hashtagUsecase.PublishHashtags(ctx, hashtags)
	if err != nil {
		panic(err)
	}

	log.Println("Send all hashtag to queue successfully")
}
