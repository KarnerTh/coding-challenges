package errors

import "fmt"

type InternalError struct {
	Msg string
	Err error
}

func (e InternalError) Error() string {
	return fmt.Sprintf("Internal error: %s - %v", e.Msg, e.Err)
}

type NotFoundError struct {
	Id string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("No item found with id %s", e.Id)
}

type ConflictError struct {
	Id string
}

func (e ConflictError) Error() string {
	return fmt.Sprintf("Item with id %s already exists", e.Id)
}

type BadInputError struct {
	Msg string
}

func (e BadInputError) Error() string {
	return e.Msg
}
