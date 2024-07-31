package sub

import "errors"

var (
	ErrInvalidEmail       = errors.New("invalid subscriber's email")
	ErrAlreadyExists      = errors.New("subscriber already exists")
	ErrNotExists          = errors.New("subscriber do not exists")
	ErrFailedToSave       = errors.New("failed to save subsriber")
	ErrFailedToUnsubscrbe = errors.New("failed to unsubscribe subsriber")
)
