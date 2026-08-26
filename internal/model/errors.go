package model

import "errors"

var (
	ErrNotFound        = errors.New("foram record not found")
	ErrConflict        = errors.New("foram record conflict")
	ErrInvalidInput    = errors.New("foram invalid input")
	ErrInvalidState    = errors.New("foram invalid state")
	ErrAlreadyReviewed = errors.New("foram section already reviewed")
	ErrChecksum        = errors.New("foram checksum mismatch")
	ErrQueueClosed     = errors.New("foram quality queue closed")
)
