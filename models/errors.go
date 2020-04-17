package models

import "strings"

const (
	ErrNotFound          modelError = "models: resource not found"
	ErrIDInvalid         modelError = "models: ID provided was invalid"
	ErrPasswordIncorrect modelError = "models: password provided was invalid"
	ErrEmailRequired     modelError = "models: email address is required"
	ErrEmailInvalid      modelError = "models: email format is not valid"
	ErrEmailTaken        modelError = "models: email address is already taken"
	ErrPasswordRequired  modelError = "models: password is required"
	ErrRememberRequired  modelError = "models: remember token is required"
	ErrPasswordToShort   modelError = "models: password must be at least 8 characters long"
	ErrRememberToShort   modelError = "models: remember token must be at least 32 bytes"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
