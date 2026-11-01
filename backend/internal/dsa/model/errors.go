package model

import "errors"

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("resource already exists")
	ErrForbidden    = errors.New("forbidden")
)
