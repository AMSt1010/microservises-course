package model

import "errors"

var (
	ErrPartNotFound      = errors.New("part not found")
	ErrUUIDCannotBeEmpty = errors.New("Uuid cannot be empty")
)
