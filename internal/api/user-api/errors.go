package uapi

import "errors"

var (
	errInvalidID    = errors.New("invalid id")
	errEmptyName    = errors.New("empty name")
	errEmptySurname = errors.New("empty surname")
)
