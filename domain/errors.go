package domain

import "errors"

var (
	ErrStreamNotFound      = errors.New("LogStream not found")
	ErrStreamAlreadyExists = errors.New("LogStream already exists")
)
