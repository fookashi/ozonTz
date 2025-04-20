package repository

import "errors"

var (
	ErrContextCanceled = errors.New("Context canceled")
	ErrNotFound        = errors.New("No records found")
	ErrPartiallyFound  = errors.New("Partially found records")
)
