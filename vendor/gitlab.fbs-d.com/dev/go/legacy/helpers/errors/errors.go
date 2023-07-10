package errors

import (
	"bufio"
	"encoding/json"
	"fmt"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
)

const (
	DetectError         = -1
	InternalError       = 1000
	ValidationError     = 1001
	BadRequest          = 1002
	RequestTimeout      = 1003
	AccessDenied        = 1004
	OperationNotAllowed = 1005
	ResourceNotFound    = 1006
	AlreadyExists       = 1007
)

var errorMessages = map[int]string{
	InternalError:       "Internal error",
	ValidationError:     "Field errors",
	BadRequest:          "Bad request",
	RequestTimeout:      "Request timeout",
	AccessDenied:        "Access denied",
	OperationNotAllowed: "Operation not allowed",
	ResourceNotFound:    "Not found",
	AlreadyExists:       "Already exists",
}

var errorStatuses = map[int]int{
	InternalError:       http.StatusInternalServerError,
	ValidationError:     http.StatusBadRequest,
	BadRequest:          http.StatusBadRequest,
	RequestTimeout:      http.StatusGatewayTimeout,
	AccessDenied:        http.StatusUnauthorized,
	OperationNotAllowed: http.StatusMethodNotAllowed,
	ResourceNotFound:    http.StatusNotFound,
	AlreadyExists:       http.StatusConflict,
}

type HttpErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code        int            `json:"code"`
	Message     string         `json:"message"`
	FieldErrors FieldErrorsOld `json:"fieldErrors,omitempty"`
}

type FieldErrorsOld []FieldErrorOld

type FieldErrorOld struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type InvalidParamError struct {
	message string
}

func (e InvalidParamError) Error() string {
	return e.message
}

func NewInvalidParamError(text string) InvalidParamError {
	var e InvalidParamError
	e.message = text
	return e
}

func NewError(err error, code int) Error {
	apiError := Error{Code: code}
	if code == DetectError {
		if errs, ok := err.(validator.ErrorMap); ok {
			apiError.Code = ValidationError
			for field, fieldError := range errs {
				apiError.FieldErrors = append(apiError.FieldErrors, FieldErrorOld{Field: field, Error: fieldError.Error()})
			}
			sort.Sort(apiError.FieldErrors)
		} else {
			apiError.Code = InternalError
		}
	}
	apiError.Message = errorMessages[apiError.Code]

	return apiError
}

func HttpApiError(err error) (status int, result HttpErrorResponse) {
	var ok bool
	var apiError Error

	if apiError, ok = err.(Error); !ok {
		apiError = NewError(err, DetectError)
	}

	status, ok = errorStatuses[apiError.Code]
	if !ok {
		status = http.StatusInternalServerError
	}

	result = HttpErrorResponse{Error: apiError}
	return
}

func (a Error) Error() string {
	return a.Message
}

func (a Error) GetCode() int {
	return a.Code
}

func (a Error) GetJson() []byte {
	res, _ := json.Marshal(a)
	return res
}

func (f FieldErrorsOld) Len() int {
	return len(f)
}

func (f FieldErrorsOld) Less(i, j int) bool {
	return f[i].Field < f[j].Field
}

func (f FieldErrorsOld) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func Oops(err error) {
	if err != nil {
		log.Println(err)
		log.Println(readErrorFileLines())
		panic(err)
	}
}

func OopsT(token string, err error) {
	if err != nil {
		log.Output(2, fmt.Sprintf("%s %v", token, err))
		os.Exit(1)
	}
}

func readErrorFileLines() map[int]string {
	_, filePath, errorLine, _ := runtime.Caller(10)
	lines := make(map[int]string)

	file, err := os.Open(filePath)
	if err != nil {
		return lines
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	currentLine := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil || currentLine > errorLine+5 {
			break
		}

		currentLine++

		if currentLine >= errorLine-5 {
			lines[currentLine] = strings.Replace(line, "\n", "", -1)
		}
	}

	return lines
}

func IsNotFoundError(err error) bool {
	if e, ok := err.(Error); ok {
		if e.Code == ResourceNotFound {
			return true
		}
	}
	return false
}
