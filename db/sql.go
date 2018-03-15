package sql

import (
	"errors"
)

var (
	ErrStatementInvalid = errors.New("sql: statement invalid")
	ErrNotAQuery        = errors.New("sql: query was not issued")
)
