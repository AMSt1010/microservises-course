package model

import "errors"

var (
	ErrPartNotFound          = errors.New("part not found")
	ErrUUIDCannotBeEmpty     = errors.New("uuid cannot be empty")
	ErrPartNameCannotBeEmpty = errors.New("part name cannot be empty")
	ErrInvalidPartPrice      = errors.New("part price must be greater than zero")
	ErrInvalidCategory       = errors.New("invalid part category")
	ErrInvalidCustomValue    = errors.New("custom value must have exactly one field populated")
)
