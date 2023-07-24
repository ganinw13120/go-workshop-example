package usecase

import (
	"context"
	"go-workshop-example/pkg/entity"
	"go-workshop-example/pkg/repository"
)

type IHashtag interface {
	GetHashtags(ctx context.Context) ([]entity.Hastag, error)
	PublishHashtags(context.Context, []entity.Hastag) error
}

type hashtag struct {
	hashtagRepo repository.IHashtag
}

func NewHashtag(hashtagRepo repository.IHashtag) *hashtag {
	return &hashtag{
		hashtagRepo: hashtagRepo,
	}
}

func (h hashtag) GetHashtags(ctx context.Context) ([]entity.Hastag, error) {
	return h.hashtagRepo.GetHashtag(ctx)
}

func (h hashtag) PublishHashtags(ctx context.Context, hashtags []entity.Hastag) error {
	for _, hashtag := range hashtags {
		err := h.hashtagRepo.PublishHashtag(ctx, hashtag)
		if err != nil {
			return err
		}
	}
	return nil
}
