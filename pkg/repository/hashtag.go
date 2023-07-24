package repository

import (
	"context"
	"encoding/json"
	"go-workshop-example/pkg/adapter"
	"go-workshop-example/pkg/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type IHashtag interface {
	GetHashtag(ctx context.Context) ([]entity.Hastag, error)
	PublishHashtag(context.Context, entity.Hastag) error
}

type hashtag struct {
	mongoDBAdapter    adapter.IMongoDBAdapter
	hashtagCollection adapter.IMongoCollection
	rabbitmqAdapter   adapter.IRabbitMQAdapter
}

func NewHashtag(mongoDBAdapter adapter.IMongoDBAdapter, hashtagCollection adapter.IMongoCollection, rabbitmqAdapter adapter.IRabbitMQAdapter) *hashtag {
	return &hashtag{
		mongoDBAdapter:    mongoDBAdapter,
		hashtagCollection: hashtagCollection,
		rabbitmqAdapter:   rabbitmqAdapter,
	}
}

func (h hashtag) GetHashtag(ctx context.Context) ([]entity.Hastag, error) {
	var hashtag []entity.Hastag
	err := h.mongoDBAdapter.Find(ctx, h.hashtagCollection, &hashtag, bson.D{{}}, nil)
	if err != nil {
		return nil, err
	}
	return hashtag, nil
}

func (h hashtag) PublishHashtag(ctx context.Context, hashtag entity.Hastag) error {
	body, err := json.Marshal(hashtag)
	if err != nil {
		return err
	}
	err = h.rabbitmqAdapter.Publish(ctx, "workshop:hashtag", "workshop:hashtag:job", body)
	return err
}
