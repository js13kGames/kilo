package token

var errInvalidInput = &InvalidInputError{}

type InvalidInputError struct{}

func (e *InvalidInputError) Error() string { return "invalid input" }
