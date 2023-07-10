package rmq

import (
	"sort"
)

type (
	msgErrors struct {
		errors []MsgError
	}

	MsgError struct {
		Index int
		Error error
	}
)

func NewMsgErrors(errors ...MsgError) (err *msgErrors) {
	var msgErrors = &msgErrors{
		errors: make([]MsgError, 0),
	}

	msgErrors.errors = append(msgErrors.errors, errors...)

	return msgErrors
}

// Error
func (m *msgErrors) Error() string {
	return "batch error"
}

// Errors возвращает ошибки, отсортированные по убыванию
func (m *msgErrors) Errors() (errors []MsgError) {
	errors = m.errors

	sort.Slice(errors, func(i, j int) bool {
		return errors[i].Index > errors[j].Index
	})

	return errors
}

// Count возвращает количество ошибок
func (m *msgErrors) Count() int {
	return len(m.errors)
}

// MsgErrorsCount
func MsgErrorsCount(err error) int {
	switch e := err.(type) {
	case *msgErrors:
		return e.Count()
	default:
		return -1
	}
}

// Add добавление ошибки
func (m *msgErrors) Add(msgErr MsgError) {
	m.errors = append(m.errors, msgErr)
}

// Add добавление ошибки
func (m *msgErrors) IsEmpty() bool {
	return len(m.errors) == 0
}
