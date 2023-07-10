package excp

import (
	"errors"
	"net/http"
)

const (
	// errors with specific httpStatus
	jsonUnmarshalInputDataException = 10100
	validationInputDataException    = 10101
	requestTimeoutException         = 10102
	multipleResultsException        = 10105
	notFoundException               = 30006
	temporaryBlockedException       = 20070
	operationNotAllowedException    = 11000
	validationException             = 31002
	countryRestrictedException      = 33000
	ipRestrictedException           = 33001
	userBlockedException            = 33002
)

func parseCodeAndHttpStatusFallback(raw map[string]interface{}) (ex IException, err error) {
	if _, ok := raw["code"]; !ok {
		err = errors.New("no code")
		return
	}

	code, err := parseInt(raw["code"])
	if err != nil {
		return
	}

	httpStatus := getHttpStatus(code)

	var message string
	if msg, ok := raw["message"]; ok {
		message = msg.(string)
	}

	if f, ok := raw["fieldErrors"]; ok {
		errorMap := ErrorMap{}
		if fieldErrors, ok := f.([]interface{}); ok {
			for _, fieldError := range fieldErrors {
				fieldErrorMap := fieldError.(map[string]interface{})
				tmsg := fieldErrorMap["error"].(string)
				errorMap[fieldErrorMap["field"].(string)] = ErrorArray{
					NewTranslatedError(errors.New(message), tmsg, tmsg),
				}
			}
		}
		return NewMapException(errorMap, code, httpStatus), nil
	}

	e := NewTranslatedError(errors.New(message), message, message)
	return NewException(e, code, httpStatus), nil
}

func getHttpStatus(code int) int {
	switch code {
	case jsonUnmarshalInputDataException,
		validationInputDataException:
		return http.StatusBadRequest
	case notFoundException:
		return http.StatusNotFound
	case requestTimeoutException:
		return http.StatusGatewayTimeout
	case multipleResultsException:
		return http.StatusMultipleChoices
	case temporaryBlockedException:
		return http.StatusConflict
	case validationException:
		return http.StatusUnprocessableEntity
	case countryRestrictedException,
		ipRestrictedException,
		userBlockedException,
		operationNotAllowedException:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
