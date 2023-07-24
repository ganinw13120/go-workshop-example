package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	Id              primitive.ObjectID `bson:"_id" json:"_id"`
	DisplayName     string             `bson:"display_name" json:"display_name"`
	Username        string             `bson:"username" json:"username"`
	ProfileImageUrl string             `bson:"profile_image_url" json:"profile_image_url"`
	Description     string             `bson:"description" json:"description"`
	Follower        int                `bson:"follower" json:"follower"`
	Following       int                `bson:"following" json:"following"`
}
