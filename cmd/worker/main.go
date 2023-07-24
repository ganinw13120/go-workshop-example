package main

import (
	"encoding/json"
	"fmt"
	"go-workshop-example/pkg/adapter"
	"go-workshop-example/pkg/entity"
	"go-workshop-example/pkg/repository"
	"go-workshop-example/pkg/usecase"
)

func main() {
	rabbitmqAdapter := adapter.NewRabbitMQAdapter()
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

		timelineUsecase.Publish(results)

		fmt.Printf("Publish %d message to queue successfully\n", len(results))
	})
}
