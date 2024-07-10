package event

import "time"

type Type string

const (
	Delete = Type("subscriber-deleted")
	Add    = Type("subscriber-added")
)

type Event struct {
	ID        int       `db:"id"`
	Type      Type      `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	Payload   string    `db:"payload"`
}
