package model

import "errors"

var (
	ErrPartNotFound      = errors.New("part not found")
	ErrUUIDCannotBeEmpty = errors.New("uuid cannot be empty")
)
