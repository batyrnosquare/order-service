package storage

import "errors"

var (
	ErrorOrderAlreadyExists = errors.New("order already exists")
	ErrorOrderNotFound      = errors.New("order not found")
)
