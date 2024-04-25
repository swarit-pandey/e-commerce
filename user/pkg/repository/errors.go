package repository

import "errors"

var (
	ErrCreate        = errors.New("failed to add entity")
	ErrRead          = errors.New("failed to get entity")
	ErrUpdate        = errors.New("failed to update entity")
	ErrDelete        = errors.New("failed to delete entity")
	ErrList          = errors.New("failed to list entity")
	ErrExists        = errors.New("failed to check if entity exists")
	ErrDoesNotExists = errors.New("entity does not exists")
)
