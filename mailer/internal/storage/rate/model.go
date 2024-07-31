package rate

import "go.mongodb.org/mongo-driver/bson/primitive"

type Exchange struct {
	ID   primitive.ObjectID `json:"_id,omitempty"`
	From string             `json:"from"`
	To   string             `json:"to"`
	Rate float32            `json:"rate"`
}
