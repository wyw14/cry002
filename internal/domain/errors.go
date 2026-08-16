package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrConflict        = errors.New("conflict")
	ErrForbidden       = errors.New("forbidden")
	ErrValidation      = errors.New("validation failed")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrTokenRevoked    = errors.New("token revoked")
	ErrTokenReplayed   = errors.New("token replayed")
	ErrInvalidState    = errors.New("invalid state transition")
	ErrAlreadyBorrowed = errors.New("case file already borrowed")
	ErrIntegrity       = errors.New("integrity check failed")
)
