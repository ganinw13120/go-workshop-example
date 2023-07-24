package entity

type Thread struct {
	Id           string  `bson:"_id" json:"_id"`
	Text         string  `bson:"text" json:"text"`
	UserId       string  `bson:"user_id" json:"user_id"`
	Likes        int     `bson:"likes" json:"likes"`
	ParentThread *string `bson:"parent_thread" json:"parent_thread"`
	RepostCount  int     `bson:"repost_count" json:"repost_count"`
}

type ThreadPayload struct {
	Thread  Thread  `json:"thread"`
	Account Account `json:"account"`
}
