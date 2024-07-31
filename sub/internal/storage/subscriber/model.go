package subscriber

import (
	"time"
)

// Subscriber is a model, which represents
// user, subscribed to daily receive mails about
// USD -> UAH rate exchanges.
type Subscriber struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}
