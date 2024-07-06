package subscriber

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscriber struct {
	ID    primitive.ObjectID `json:"_id,omitempty"`
	Email string             `json:"email"`
}
