package excp

import (
	"fmt"
)

// Map ошибок (внутри массивы ошибок), реализует интерфейс error
type ErrorMap map[string]ErrorArray

func (err ErrorMap) Error() string {
	for k, errs := range err {
		if len(errs) > 0 {
			return fmt.Sprintf("%s: %s", k, errs.Error())
		}
	}

	return ""
}

// Массив ошибок, реализует интерфейс error
type ErrorArray []ITranslatedError

func (err ErrorArray) Error() string {
	if len(err) > 0 {
		return err[0].Error()
	}
	return ""
}

type translatedError struct {
	err      error
	message  string
	tMessage string
}

type ITranslatedError interface {
	Error() string
	Message() string
	getError() error
	getMap() map[string]interface{}
}

func NewTranslatedError(err error, message, tMessage string) ITranslatedError {
	return translatedError{
		err:      err,
		message:  message,
		tMessage: tMessage,
	}
}

func (e translatedError) Error() string {
	return e.message
}

func (e translatedError) Message() string {
	return e.tMessage
}

func (e translatedError) getHttpJson() []byte {
	return []byte("")
}

func (e translatedError) getError() error {
	return e.err
}

func (e translatedError) getMap() map[string]interface{} {
	errorMessage := ""
	if e.err != nil {
		errorMessage = e.err.Error()
	}
	return map[string]interface{}{
		"error":    errorMessage,
		"message":  e.message,
		"tMessage": e.tMessage,
	}
}

type ITranslatableError interface {
	error
	Params() map[string]interface{}
	Category() string
}

type translatableError struct {
	category string
	message  string
	params   map[string]interface{}
}

func (e translatableError) Error() string {
	return e.message
}

func (e translatableError) Params() map[string]interface{} {
	return e.params
}

func (e translatableError) Category() string {
	return e.category
}

func NewError(message string, params map[string]interface{}) ITranslatableError {
	return translatableError{
		message: message,
		params:  params,
	}
}

func NewErrorV2(category string, message string, params map[string]interface{}) ITranslatableError {
	return translatableError{
		category: category,
		message:  message,
		params:   params,
	}
}
