package toolkit

import (
	"errors"
	"gitlab.fbs-d.com/dev/go/legacy/exceptions"
	"gitlab.fbs-d.com/dev/go/legacy/i18n"
	"net/http"
)

const (
	NewGenerateRandomStringExceptionCode = 10108
	NewNilDataExceptionCode              = 10109
	NewUuidExceptionCode                 = 10110
	CurrencyConversionException          = 11001
	DataTypeConversionExceptionCode      = 10111
	ParseTimeException                   = 11100
	DecodeImageExceptionCode             = 11101
	SaveImageExceptionCode               = 11102
)

func NewGenerateRandomStringException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.generate-random-string-error", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		NewGenerateRandomStringExceptionCode,
		http.StatusInternalServerError,
	)
}

func NewNilDataException() excp.IException {
	t := i18n.T(i18nCategory, "error.nil-data-error", "")
	return excp.NewException(
		excp.NewTranslatedError(errors.New(""), t, t),
		NewNilDataExceptionCode,
		http.StatusInternalServerError,
	)
}

func NewUuidException() excp.IException {
	t := i18n.T(i18nCategory, "error.generate-uuid-error", "")
	return excp.NewException(
		excp.NewTranslatedError(errors.New(""), t, t),
		NewUuidExceptionCode,
		http.StatusInternalServerError,
	)
}

func NewCurrencyConversionException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.currency-conversion", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		CurrencyConversionException,
		http.StatusInternalServerError,
	)
}

func NewParseTimeException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.parse-time", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		ParseTimeException,
		http.StatusInternalServerError,
	)
}

func NewDataTypeConversionException() excp.IException {
	t := i18n.T(i18nCategory, "error.data-type-conversion", "")
	return excp.NewException(
		excp.NewTranslatedError(errors.New(""), t, t),
		DataTypeConversionExceptionCode,
		http.StatusInternalServerError,
	)
}

func NewDecodeImageException(err error, l string) excp.IException {
	t := i18n.T(i18nCategory, "error.decode-image", l)
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		DecodeImageExceptionCode,
		http.StatusBadRequest,
	)
}

func NewSaveImageException(err error, l string) excp.IException {
	t := i18n.T(i18nCategory, "error.save-image", l)
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		SaveImageExceptionCode,
		http.StatusInternalServerError,
	)
}
