package domain

import "errors"

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountNotFound   = errors.New("account not found")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrSameAccount       = errors.New("cannot transfer to the same account")
)
