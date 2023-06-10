package errorsapp

import "errors"

var ErrLoginAlreadyBorrowed = errors.New("login already busy")

var ErrNoFoundUser = errors.New("no user found")
