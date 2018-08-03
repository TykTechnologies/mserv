package errors

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

func New(s string) error {
	return errors.New(s)
}
