package model

import (
	"errors"
)

var (
	ErrAlreadyExists                = errors.New("order already exists")
	ErrOrderNotFound                = errors.New("order not found")
	ErrOrderUUIDCannotBeEmpty       = errors.New("order uuid can not be empty")
	ErrUserUUIDCannotBeEmpty        = errors.New("user uuid can not be empty")
	ErrPartsUUIDsCannotBeEmpty      = errors.New("parts uuids cannot be empty")
	ErrAtLeastOnePartWasNotFound    = errors.New("at least one part was not found")
	ErrInventoryUnavailable         = errors.New("inventory service unavailable")
	ErrPaymentUnavailable           = errors.New("payment service unavailable")
	ErrTransactionUUIDCannotBeEmpty = errors.New("transaction uuid cannot be empty")
	ErrInsufficientPaymentMethod    = errors.New("insufficient payment method")
	ErrAlreadyPaid                  = errors.New("order already paid")
	ErrOrderCancelled               = errors.New("order already cancelled")
	ErrCustomValueExactOneField     = errors.New("custom value must have exactly one field filled")
)
