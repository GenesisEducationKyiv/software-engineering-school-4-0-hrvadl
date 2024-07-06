package rate

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrFailedToReplace = errors.New("failed to replace")
	ErrFailetToGet     = errors.New("failed to get")
)
