package service

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUnauthorized = errors.New("unauthorized")
var ErrInvalidInput = errors.New("invalid input")
var ErrNonexistentClient = errors.New("client does not exisit")