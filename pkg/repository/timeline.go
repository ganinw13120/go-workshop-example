package repository

import (
	"encoding/json"
	"fmt"
	"go-workshop-example/pkg/adapter"
	"go-workshop-example/pkg/entity"
	"io/ioutil"
	"net/http"
)

const PageSize = 100

type ITimeline interface {
	FetchTimeline(hashtag string, cursor *string) ([]entity.Thread, *string, error)
}

type timeline struct {
	rabbitmqAdapter adapter.MessageBroker
}

func NewTimeline(rabbitmqAdapter adapter.MessageBroker) *timeline {
	return &timeline{
		rabbitmqAdapter: rabbitmqAdapter,
	}
}

type TimelineResponse struct {
	Data     []entity.Thread `json:"data"`
	NextPage *string         `json:"next_page"`
}

// Todo -> Go to adapter
func (timeline) FetchTimeline(hashtag string, cursor *string) ([]entity.Thread, *string, error) {
	url := fmt.Sprintf("https://go-workshop-2zcpzmfnyq-de.a.run.app/thread/?hashtag=%s&page_size=%d", hashtag, PageSize)
	if cursor != nil {
		url += fmt.Sprintf("&cursor=%s", *cursor)
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	var response TimelineResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, nil, err
	}
	return response.Data, response.NextPage, nil
}
