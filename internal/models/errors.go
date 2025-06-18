package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: No mathching record found!")

	ErrInvalidCredential = errors.New("models: invalid creadentials!")

	ErrDuplicateEmail = errors.New("models: duplicate emails")
)
