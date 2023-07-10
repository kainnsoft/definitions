package excp

import "encoding/json"

// реализует IException
type mapException struct {
	errors     ErrorMap
	code       int
	httpStatus int
}

func (ex mapException) Error() string {
	return ex.errors.Error()
}

func (ex mapException) Code() int {
	return ex.code
}

func (ex mapException) HttpStatus() int {
	return ex.httpStatus
}

func (ex mapException) GetJson() []byte {
	errMap := make(map[string][]interface{})
	for field, errs := range ex.errors {
		errMap[field] = make([]interface{}, len(errs))
		for i, err := range errs {
			errMap[field][i] = err.getMap()
		}
	}

	data := map[string]map[string]interface{}{
		"mapException": {
			"errors":     errMap,
			"code":       ex.code,
			"httpStatus": ex.httpStatus,
		},
	}
	r, _ := json.Marshal(data)
	return r
}

func (ex mapException) GetHttpJson() []byte {
	errMap := make(map[string][]string)
	for field, errs := range ex.errors {
		errMap[field] = make([]string, len(errs))
		for i, err := range errs {
			errMap[field][i] = err.Message()
		}
	}

	data := map[string]map[string]interface{}{
		"error": {
			"code":   ex.code,
			"errors": errMap,
		},
	}
	r, _ := json.Marshal(data)
	return r
}

func (ex mapException) GetError() error {
	return ex.errors
}

func (ex mapException) GetTranslatedMessage() string {
	for _, e := range ex.errors {
		return e[0].Message()
	}
	return ""
}

func NewMapException(errors ErrorMap, code, httpStatus int) IException {
	return mapException{
		errors:     errors,
		code:       code,
		httpStatus: httpStatus,
	}
}
