package subscriber

import "errors"

var (
	ErrEmtpyEmail     = errors.New("email can not be empty")
	ErrFailedToGet    = errors.New("failed to get")
	ErrFailedToSave   = errors.New("failed to save")
	ErrFailedToDelete = errors.New("failed to delete")
	ErrNoSubscribers  = errors.New("no subscribers")
)
