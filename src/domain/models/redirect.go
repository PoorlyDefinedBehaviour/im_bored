package models

type Redirect struct {
	Code      string `json:"code" bson:"code" msgpack:"code"`
	Url       string `json:"url" bson:"url" msgpack:"url" validate:"empty=false & format=url"`
	CreatedAt int64  `json:"createdAt" bson:"createdAt" msgpack:"createdAt"`
}
