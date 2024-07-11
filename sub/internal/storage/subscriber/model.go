package subscriber

import (
	"time"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
)

func New(mail string) Subscriber {
	return Subscriber{
		Email: mail,
		event: &event.Event{
			Type:    event.Added,
			Payload: mail,
		},
	}
}

// Subscriber is a model, which represents
// user, subscribed to daily receive mails about
// USD -> UAH rate exchanges.
type Subscriber struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	event     *event.Event
}

func (s Subscriber) GetEvent() event.Event {
	if s.event != nil {
		return *s.event
	}

	return event.Event{
		Type: event.Deleted, Payload: s.Email,
	}
}
