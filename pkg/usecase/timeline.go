package usecase

import (
	"fmt"
	"go-workshop-example/pkg/entity"
	"go-workshop-example/pkg/repository"
	"sync"
)

type ITimeline interface {
	FetchTimelineFromHashtag(hashtag entity.Hastag) ([]*entity.ThreadPayload, error)
}

type timeline struct {
	timelineRepo repository.ITimeline
	accountRepo  repository.IAccount
}

func NewTimeline(timelineRepo repository.ITimeline, accountRepo repository.IAccount) *timeline {
	return &timeline{
		timelineRepo: timelineRepo,
		accountRepo:  accountRepo,
	}
}

func (t timeline) fetchThreadAccount(thread entity.Thread) (*entity.ThreadPayload, error) {
	account, err := t.accountRepo.FetchAccount(thread.UserId)
	if err != nil {
		return nil, err
	}
	result := entity.ThreadPayload{
		Thread:  thread,
		Account: *account,
	}
	return &result, nil
}

func (t timeline) FetchTimelineFromHashtag(hashtag entity.Hastag) ([]*entity.ThreadPayload, error) {
	var results []*entity.ThreadPayload
	var cursor *string
	var mutex sync.Mutex
	var wg sync.WaitGroup
	for {
		threads, next, err := t.timelineRepo.FetchTimeline(hashtag.Keyword, cursor)
		if err != nil {
			return nil, err
		}
		cursor = next
		for _, thread := range threads {
			go func(thread entity.Thread) {
				wg.Add(1)
				defer wg.Done()

				result, err := t.fetchThreadAccount(thread)
				if err != nil {
					fmt.Println(err)
					return
				}

				mutex.Lock()
				results = append(results, result)
				mutex.Unlock()
			}(thread)
		}
		if cursor == nil {
			break
		}
	}
	wg.Wait()
	return results, nil
}
