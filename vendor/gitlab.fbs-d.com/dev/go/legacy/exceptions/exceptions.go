package excp

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"gitlab.fbs-d.com/dev/go/legacy/i18n"
)

const (
	// JSON processing exceptions
	JsonUnmarshalInputDataException = 10100
	JsonUnmarshalException          = 10104
	JsonMarshalException            = 10105

	OperationNotAllowedException = 11000

	TemporaryBlockedException = 20070

	EventProcessException = 20100

	// Redis Exceptions
	RedisGetException    = 10050
	RedisDeleteException = 10051
	RedisSetException    = 10052
	RedisRenameException = 10053

	i18nCategory = "excp"
)

func init() {
	err := i18n.LoadRaw(messages, i18nCategory, "")
	if err != nil {
		log.Output(2, fmt.Sprintf("%s %v", "sys", err))
		os.Exit(1)
	}
}

// Общий интерфейс ошибки
type IException interface {
	error
	Code() int
	HttpStatus() int
	GetJson() []byte
	GetHttpJson() []byte
	GetError() error
	GetTranslatedMessage() string
}

func Parse(raw []byte) IException {
	res := make(map[string]map[string]interface{})
	err := json.Unmarshal(raw, &res)
	if err != nil {
		return nil
	}
	if data, ok := res["exception"]; ok {
		code, httpStatus, err := parseCodeAndHttpStatus(data)
		if err != nil {
			return nil
		}

		var tErr ITranslatedError
		if t, ok := data["error"]; ok {
			tErr = ParseTranslatedError(t)
		}

		return NewException(tErr, code, httpStatus)
	}
	if data, ok := res["mapException"]; ok {
		code, httpStatus, err := parseCodeAndHttpStatus(data)
		if err != nil {
			return nil
		}

		eMap := make(ErrorMap)
		if _, ok := data["errors"]; ok {
			if m, ok := data["errors"].(map[string]interface{}); ok {
				for k, v := range m {
					if va, ok := v.([]interface{}); ok {
						errArray := make(ErrorArray, len(va))
						for index, vi := range va {
							errArray[index] = ParseTranslatedError(vi)
						}
						eMap[k] = errArray
					}
				}
			}
		}

		return NewMapException(eMap, code, httpStatus)
	}
	if data, ok := res["error"]; ok {
		ex, err := parseCodeAndHttpStatusFallback(data)
		if err != nil {
			return nil
		}
		return ex
	}
	return nil
}

func parseCodeAndHttpStatus(raw map[string]interface{}) (code, httpStatus int, err error) {
	if _, ok := raw["code"]; !ok {
		err = errors.New("No code")
		return
	}
	if _, ok := raw["httpStatus"]; !ok {
		err = errors.New("No httpStatus")
		return
	}
	code, err = parseInt(raw["code"])
	if err != nil {
		return
	}
	httpStatus, err = parseInt(raw["httpStatus"])
	return
}

func parseInt(raw interface{}) (code int, err error) {
	if f, ok := raw.(float64); ok {
		code = int(f)
		return
	}
	err = errors.New("Failed parse float64")
	return
}

func ParseTranslatedError(raw interface{}) ITranslatedError {
	if errMap, ok := raw.(map[string]interface{}); ok {
		var m, t string
		var err error
		if i, ok := errMap["error"]; ok {
			if e, ok := i.(string); ok {
				err = errors.New(e)
			}
		}
		if i, ok := errMap["message"]; ok {
			m, _ = i.(string)
		}
		if i, ok := errMap["tMessage"]; ok {
			t, _ = i.(string)
		}
		return NewTranslatedError(err, m, t)
	}
	return nil
}

func NewJsonUnmarshalInputDataException(err error) IException {
	t := i18n.T(i18nCategory, "error.json-unmarshal-input-data", "")
	return NewException(NewTranslatedError(err, t, t), JsonUnmarshalInputDataException, http.StatusBadRequest)
}

func NewJsonUnmarshalException(err error) IException {
	t := i18n.T(i18nCategory, "error.json-unmarshal", "")
	return NewException(NewTranslatedError(err, t, t), JsonUnmarshalException, http.StatusInternalServerError)
}

func NewJsonMarshalException(err error) IException {
	t := i18n.T(i18nCategory, "error.json-marshal", "")
	return NewException(NewTranslatedError(err, t, t), JsonMarshalException, http.StatusInternalServerError)
}

func NewOperationNotAllowedException(err error) IException {
	t := i18n.T(i18nCategory, "error.operation-not-allowed", "")
	return NewException(NewTranslatedError(err, t, t), OperationNotAllowedException, http.StatusForbidden)
}

func NewTemporaryBlockedException(err error) IException {
	t := i18n.T(i18nCategory, "error.temporary-blocked", "")
	return NewException(NewTranslatedError(err, t, t), TemporaryBlockedException, http.StatusConflict)
}

func NewEventProcessException(err error) IException {
	t := i18n.T(i18nCategory, "error.event", "")
	return NewException(NewTranslatedError(err, t, t), EventProcessException, http.StatusInternalServerError)
}

func NewRedisGetException(err error) IException {
	t := i18n.T(i18nCategory, "error.redis-get", "")
	return NewException(NewTranslatedError(err, t, t), RedisGetException, http.StatusInternalServerError)
}

func NewRedisSetException(err error) IException {
	t := i18n.T(i18nCategory, "error.redis-set", "")
	return NewException(NewTranslatedError(err, t, t), RedisSetException, http.StatusInternalServerError)
}

func NewRedisDeleteException(err error) IException {
	t := i18n.T(i18nCategory, "error.redis-delete", "")
	return NewException(NewTranslatedError(err, t, t), RedisDeleteException, http.StatusInternalServerError)
}

func NewRedisRenameException(err error) IException {
	t := i18n.T(i18nCategory, "error.redis-rename", "")
	return NewException(NewTranslatedError(err, t, t), RedisRenameException, http.StatusInternalServerError)
}
