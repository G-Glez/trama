package core

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrDB             = errors.New("database error")
	ErrDataCorruption = errors.New("data corruption")
)
