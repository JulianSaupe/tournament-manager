package domain

import (
	"errors"
)

// Error interfaces for different error types

// ErrNotFound signals that the requested object doesn't exist
type ErrNotFound interface {
	error
	NotFound()
}

// ErrInvalidParameter signals that the user input is invalid
type ErrInvalidParameter interface {
	error
	InvalidParameter()
}

// ErrUnauthorized is used to signify that the user is not authorized to perform a specific action
type ErrUnauthorized interface {
	error
	Unauthorized()
}

// ErrForbidden signals that the requested action cannot be performed under any circumstances
type ErrForbidden interface {
	error
	Forbidden()
}

type ErrNotAllowed interface {
	error
	NotAllowed()
}

// Error implementations

type errNotFound struct{ error }

func (errNotFound) NotFound()       {}
func (e errNotFound) Unwrap() error { return e.error }

type errInvalidParameter struct{ error }

func (errInvalidParameter) InvalidParameter() {}
func (e errInvalidParameter) Unwrap() error   { return e.error }

type errUnauthorized struct{ error }

func (errUnauthorized) Unauthorized()   {}
func (e errUnauthorized) Unwrap() error { return e.error }

type errForbidden struct{ error }

func (errForbidden) Forbidden()      {}
func (e errForbidden) Unwrap() error { return e.error }

type errNotAllowed struct{ error }

func (errNotAllowed) NotAllowed()     {}
func (e errNotAllowed) Unwrap() error { return e.error }

// Helper functions to create errors

// NewNotFoundError creates a new ErrNotFound from the given error or message
func NewNotFoundError(msg string) error {
	return errNotFound{errors.New(msg)}
}

// NewInvalidParameterError creates a new ErrInvalidParameter from the given error or message
func NewInvalidParameterError(msg string) error {
	return errInvalidParameter{errors.New(msg)}
}

// NewUnauthorizedError creates a new ErrUnauthorized from the given error or message
func NewUnauthorizedError(msg string) error {
	return errUnauthorized{errors.New(msg)}
}

// NewForbiddenError creates a new ErrForbidden from the given error or message
func NewForbiddenError(msg string) error {
	return errForbidden{errors.New(msg)}
}

// NewNotAllowedError creates a new ErrNotAllowed from the given error or message
func NewNotAllowedError(msg string) error { return errNotAllowed{errors.New(msg)} }

// Helper functions to check error types

// IsNotFound returns true if the error is an ErrNotFound
func IsNotFound(err error) bool {
	var notFound ErrNotFound
	return errors.As(err, &notFound)
}

// IsInvalidParameter returns true if the error is an ErrInvalidParameter
func IsInvalidParameter(err error) bool {
	var invalidParam ErrInvalidParameter
	return errors.As(err, &invalidParam)
}

// IsUnauthorized returns true if the error is an ErrUnauthorized
func IsUnauthorized(err error) bool {
	var unauthorized ErrUnauthorized
	return errors.As(err, &unauthorized)
}

// IsForbidden returns true if the error is an ErrForbidden
func IsForbidden(err error) bool {
	var forbidden ErrForbidden
	return errors.As(err, &forbidden)
}

func IsNotAllowed(err error) bool {
	var notAllowed ErrNotAllowed
	return errors.As(err, &notAllowed)
}
