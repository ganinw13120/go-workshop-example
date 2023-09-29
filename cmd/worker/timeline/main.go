package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go-workshop-example/pkg/adapter"
	"go-workshop-example/pkg/entity"
	"go-workshop-example/pkg/repository"
	"go-workshop-example/pkg/usecase"
	"time"
)

func main() {
	rabbitmqAdapter, err := adapter.NewRabbitMQ("", adapter.QueueConfig{})
	if err != nil {
		panic(err)
	}
	defer rabbitmqAdapter.CleanUp()

	timelineRepo := repository.NewTimeline(rabbitmqAdapter)
	accountRepo := repository.NewAccount()

	timelineUsecase := usecase.NewTimeline(timelineRepo, accountRepo)

	rabbitmqAdapter.Consume("workshop:hashtag:job", func(job []byte) {
		var hashtag entity.Hastag
		err := json.Unmarshal(job, &hashtag)
		if err != nil {
			fmt.Println(err)
			return
		}

		results, err := timelineUsecase.FetchTimelineFromHashtag(hashtag)
		if err != nil {
			fmt.Println(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		for _, result := range results {
			body, err := json.Marshal(*result)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = rabbitmqAdapter.Publish(ctx, "", "", body)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Printf("Publish %d message to queue successfully\n", len(results))
	})
}
