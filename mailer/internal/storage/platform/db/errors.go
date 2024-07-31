package db

import "errors"

var (
	ErrFailedConnect = errors.New("failed to connect")
	ErrNotFound      = errors.New("not found")
)
