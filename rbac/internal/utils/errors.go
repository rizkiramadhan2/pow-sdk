package utils

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrUnauthorized   = errors.New("unauthorized")
)
