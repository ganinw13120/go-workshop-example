package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hastag struct {
	Id      primitive.ObjectID `bson:"_id" json:"_id"`
	Keyword string             `bson:"keyword" json:"keyword"`
}
