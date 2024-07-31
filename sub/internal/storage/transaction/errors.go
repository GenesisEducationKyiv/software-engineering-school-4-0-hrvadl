package transaction

import "errors"

var (
	ErrFailedGetTx     = errors.New("failed to get tx from ctx")
	ErrFailedBeginTx   = errors.New("failed to begin tx")
	ErrFailedExecuteTx = errors.New("failed to execute tx")
)
