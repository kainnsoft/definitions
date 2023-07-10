package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type codeError struct {
	message     string
	translation string
	code        int
	httpCode    int
	cause       error
}

// New constructor implementation
func New(message string) error {
	return &codeError{
		message: message,
	}
}

// NewWithCode constructor implementation
func NewWithCode(message string, code int) error {
	return &codeError{
		code:    code,
		message: message,
	}
}

// NewWithHTTPCode constructor implementation
func NewWithHTTPCode(message string, httpCode int) error {
	return &codeError{
		httpCode: httpCode,
		message:  message,
	}
}

// NewWithCodes constructor implementation
func NewWithCodes(message string, code, httpCode int) error {
	return &codeError{
		code:     code,
		httpCode: httpCode,
		message:  message,
	}
}

// NewWithTranslations constructor implementation
func NewWithTranslations(message string, code, httpCode int, translation string) error {
	return &codeError{
		code:        code,
		httpCode:    httpCode,
		message:     message,
		translation: translation,
	}
}

//NewWithTranslation error с переводами старых/новых сервисов
func NewWithTranslation(message string, code, httpCode int, oldTranslation, newTranslation string) error {
	var translation string
	translation = oldTranslation
	if translation == "" {
		translation = newTranslation
	}
	if translation == "" {
		translation = message
	}

	return &codeError{
		code:        code,
		httpCode:    httpCode,
		message:     message,
		translation: translation,
	}
}

// Wrap wrap error with message
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &codeError{
		cause:   err,
		message: message,
	}
}

// WrapCode wrap error with custom code
func WrapCode(code int, err error) error {
	return &codeError{
		code:  code,
		cause: err,
	}
}

// WrapHTTPCode wrap error with http code
func WrapHTTPCode(httpCode int, err error) error {
	return &codeError{
		httpCode: httpCode,
		cause:    err,
	}
}

// WrapCodes wrap error with custom and http codes
func WrapCodes(code, httpCode int, err error) error {
	return &codeError{
		code:     code,
		httpCode: httpCode,
		cause:    err,
	}
}

// WrapTranslation wrap error with translation
func WrapTranslation(code, httpCode int, err error) error {
	return &codeError{
		code:     code,
		httpCode: httpCode,
		cause:    err,
	}
}

func WrapCause(err *codeError, cause error) error {
	return &codeError{
		code:        err.code,
		httpCode:    err.httpCode,
		message:     err.message,
		translation: err.translation,
		cause:       err,
	}
}

// Code getter
func (c *codeError) Code() int {
	return c.code
}

// Translation getter
func (c *codeError) Translation() string {
	if c.translation == "" {
		return c.message
	}
	return c.translation
}

// Error getter
func (c *codeError) Error() string {
	if c.cause == nil {
		return c.message
	}
	return c.message + ": " + c.cause.Error()
}
func (c *codeError) Cause() error { return c.cause }

// HTTPCode getter
func (c *codeError) HTTPCode() int {
	return c.httpCode
}

// Code extractor
func Code(err error) int {
	type codeKeeper interface {
		Code() int
	}

	code, ok := err.(codeKeeper)
	if !ok {
		return -1
	}
	return code.Code()
}

// HTTPCode extractor
func HTTPCode(err error) int {
	type codeKeeper interface {
		HTTPCode() int
	}

	code, ok := err.(codeKeeper)
	if !ok {
		return -1
	}
	return code.HTTPCode()
}

// Translation extractor
func Translation(err error) string {
	type codeKeeper interface {
		Translation() string
	}

	code, ok := err.(codeKeeper)
	if !ok {
		return ""
	}
	return code.Translation()
}

func GRPCError(err, defaultErr error) error {
	var code int
	if code = Code(err); code > 0 {
		var message = err.Error()

		if t := Translation(err); len(t) > 0 {
			message = t
		}

		return status.Error(codes.Code(code), message)
	}

	return defaultErr
}
