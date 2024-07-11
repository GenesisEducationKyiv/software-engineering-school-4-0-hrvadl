package event

import "time"

type Type string

const (
	Deleted = Type("subscriber-deleted")
	Added   = Type("subscriber-added")
)

type Event struct {
	ID        int       `db:"id"`
	Type      Type      `db:"type"`
	CreatedAt time.Time `db:"created_at"`
	Payload   string    `db:"payload"`
}
