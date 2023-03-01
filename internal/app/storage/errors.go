package storage

import "errors"

// Typical errors
var (
	ErrShortURLRemoved = errors.New("short URL has been removed")
	ErrAlreadyExists   = errors.New("such an entry exists in the database")
	ErrDBConnection    = errors.New("failed to check database connection")
	ErrInvalidData     = errors.New("data invalid")
	ErrUIDNotValid     = errors.New("short URL does not exist")
	ErrUserIDNotValid  = errors.New("invalid user id")
	ErrUserNotExists   = errors.New("user with current id does not exist")
	ErrUserUrlsEmpty   = errors.New("shortened URLs for current user is empty")
)
