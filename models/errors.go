package models

import "strings"

const (
	ErrNotFound          modelError = "models: resource not found"
	ErrPasswordIncorrect modelError = "models: password provided was invalid"
	ErrEmailRequired     modelError = "models: email address is required"
	ErrEmailInvalid      modelError = "models: email format is not valid"
	ErrEmailTaken        modelError = "models: email address is already taken"
	ErrPasswordRequired  modelError = "models: password is required"
	ErrPasswordToShort   modelError = "models: password must be at least 8 characters long"
	ErrTitleRequired     modelError = "models: title is required"

	ErrIDInvalid        privateError = "models: ID provided was invalid"
	ErrRememberRequired privateError = "models: remember token is required"
	ErrRememberToShort  privateError = "models: remember token must be at least 32 bytes"
	ErrUserIDRequired   privateError = "models: user ID is required"
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

type privateError string

func (e privateError) Error() string {
	return string(e)
}
