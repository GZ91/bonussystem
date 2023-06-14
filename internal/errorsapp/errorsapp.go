package errorsapp

import "errors"

var ErrLoginAlreadyBorrowed = errors.New("login already busy")

var ErrNoFoundUser = errors.New("no user found")

var ErrIncorrectOrderNumber = errors.New("incorrect order number")

var ErrOrderAlreadyThisUser = errors.New("this order has already been entered by this user")

var ErrOrderAlreadyAnotherUser = errors.New("this order has already been entered by another user")

var ErrNoRecords = errors.New("no records")
