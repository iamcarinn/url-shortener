package storage

import "errors"

// Общие ошибки
var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

type Storage interface {
	SaveURL(urlToSave string, alias string) (error)
	GetURL(alias string) (string, error)
}
