package excp

import "encoding/json"

// реализует IException
type exception struct {
	err        ITranslatedError
	code       int
	httpStatus int
}

func (ex exception) Error() string {
	if ex.err != nil {
		err := ex.GetError()
		if err != nil {
			return err.Error()
		}
		return "unknown error"
	}
	return ""
}

func (ex exception) Code() int {
	return ex.code
}

func (ex exception) HttpStatus() int {
	return ex.httpStatus
}

func (ex exception) GetJson() []byte {
	data := map[string]map[string]interface{}{
		"exception": {
			"error":      ex.err.getMap(),
			"code":       ex.code,
			"httpStatus": ex.httpStatus,
		},
	}
	r, _ := json.Marshal(data)
	return r
}

func (ex exception) GetHttpJson() []byte {
	data := map[string]map[string]interface{}{
		"error": {
			"code":    ex.code,
			"message": ex.err.Message(),
		},
	}
	r, _ := json.Marshal(data)
	return r
}

func (ex exception) GetError() error {
	return ex.err.getError()
}

func (ex exception) GetTranslatedMessage() string {
	return ex.err.Message()
}

func NewException(err ITranslatedError, code, httpStatus int) IException {
	return exception{
		err:        err,
		code:       code,
		httpStatus: httpStatus,
	}
}
