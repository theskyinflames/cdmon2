package hosting

import "errors"

var (
	DbErrorAlreadyExist error = errors.New("already exist")
	DbErrorNotFound     error = errors.New("not found")
)
