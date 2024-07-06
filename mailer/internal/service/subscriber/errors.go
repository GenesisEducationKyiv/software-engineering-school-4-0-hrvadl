package subscriber

import "errors"

var (
	ErrFailedToGet    = errors.New("failed to get")
	ErrFailedToSave   = errors.New("failed to save")
	ErrFailedToDelete = errors.New("failed to delete")
	ErrNoSubscribers  = errors.New("no subscribers")
)
