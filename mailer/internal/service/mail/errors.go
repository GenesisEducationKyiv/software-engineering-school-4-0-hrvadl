package mail

import "errors"

var (
	ErrEmptyContent   = errors.New("content can not be empty")
	ErrEmptySubject   = errors.New("subject can not be empty")
	ErrEmptyReceivers = errors.New("need to specify at least 1 receiver")
)
