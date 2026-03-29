package domain

import "fmt"

type ErrorKind int

const (
	ErrNotFound ErrorKind = iota
	ErrValidation
	ErrConflict
)

type Error struct {
	Kind    ErrorKind
	Message string
	Err     error
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Err }

func NotFoundErr(entity, id string) *Error {
	return &Error{Kind: ErrNotFound, Message: fmt.Sprintf("%s %s not found", entity, id)}
}

func ValidationErr(msg string) *Error {
	return &Error{Kind: ErrValidation, Message: msg}
}

func ConflictErr(msg string) *Error {
	return &Error{Kind: ErrConflict, Message: msg}
}
