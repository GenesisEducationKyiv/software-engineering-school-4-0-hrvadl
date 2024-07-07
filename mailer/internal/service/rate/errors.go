package rate

import "errors"

var (
	ErrEmptyRate       = errors.New("rate is empty")
	ErrNotFound        = errors.New("not found")
	ErrFailedToReplace = errors.New("failed to replace")
	ErrFailetToGet     = errors.New("failed to get")
)
