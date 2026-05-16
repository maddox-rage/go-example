package storage 

import "errors"

var (
	ErrURLFound = errors.New("url not found")
	ErrURLExist = errors.New("url exist")
)
