package domain

import "errors"

var (
	ErrStreamNotFound      = errors.New("LogStream not found")
	ErrClientNotFound      = errors.New("Client not found")
	ErrStreamAlreadyExists = errors.New("LogStream already exists")
)
