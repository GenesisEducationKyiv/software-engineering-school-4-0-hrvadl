package subscriber

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/event"
	"github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/storage/platform/db"
)

const (
	deleteEvent = "delete-subscriber"
	insertEvent = "add-subscriber"
)

// NewRepo constructs repo with provided sqlx DB connection.
// NOTE: it expects db connection to be connection MySQL.
func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

// Repo is a thin abstraction to not do sqlx queries
// directly in the services. Therefore specific underlying DB could
// be more easily changed in the future.
type Repo struct {
	db *sqlx.DB
}

// Save method saves subscriber to the repo and then returns
// newly created ID. Could return an error if email is not valid, or such email
// already exists.
func (r *Repo) Save(ctx context.Context, s Subscriber) (int64, error) {
	const query = "INSERT INTO subscribers (email) VALUES (?)"

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = errors.Join(err, tx.Commit())
		}
	}()

	res, err := tx.ExecContext(ctx, query, s.Email)
	if err != nil {
		var mySQLErr *mysql.MySQLError
		if errors.As(err, &mySQLErr) && mySQLErr.Number == db.AlreadyExistsErrCode {
			return 0, ErrAlreadyExists
		}

		return 0, err
	}

	e := event.Event{Type: insertEvent, Payload: s.Email}
	es := event.NewRepo(tx)
	if exist, err := es.GetByPayload(ctx, deleteEvent, s.Email); err == nil && exist != nil {
		err := es.DeleteByID(ctx, exist.ID)
		id, idErr := res.LastInsertId()
		return id, errors.Join(err, idErr)
	}

	if err := es.Save(ctx, e); err != nil {
		return 0, fmt.Errorf("failed to save event: %w", err)
	}

	return res.LastInsertId()
}

// Get method gets all subscribers from the DB.
func (r *Repo) Get(ctx context.Context) ([]Subscriber, error) {
	subscribers := []Subscriber{}
	if err := r.db.SelectContext(ctx, &subscribers, "SELECT * FROM subscribers"); err != nil {
		return nil, err
	}

	return subscribers, nil
}

// GetByEmail method gets subscriber from the DB by his email.
func (r *Repo) GetByEmail(ctx context.Context, email string) (*Subscriber, error) {
	const query = "SELECT * FROM subscribers WHERE email = ? LIMIT 1"
	subscriber := Subscriber{}
	if err := r.db.GetContext(ctx, &subscriber, query, email); err != nil {
		return nil, err
	}

	return &subscriber, nil
}

// DeleteByEmail method gets subscriber from the DB by his email.
func (r *Repo) DeleteByEmail(ctx context.Context, email string) error {
	const query = "DELETE FROM subscribers WHERE email = ?"

	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf(" failed to create transaction: %w", err)
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
		} else {
			err = errors.Join(err, tx.Commit())
		}
	}()

	if _, err := tx.ExecContext(ctx, query, email); err != nil {
		return fmt.Errorf("failed to delete sub: %w", err)
	}

	e := event.Event{Type: deleteEvent, Payload: email}
	es := event.NewRepo(tx)
	if exist, err := es.GetByPayload(ctx, insertEvent, email); err == nil && exist != nil {
		return errors.Join(err, es.DeleteByID(ctx, exist.ID))
	}

	if err := es.Save(ctx, e); err != nil {
		return fmt.Errorf("failed to save event: %w", err)
	}

	return nil
}
