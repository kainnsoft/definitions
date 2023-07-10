package rmq

import (
	"errors"
	"gitlab.fbs-d.com/dev/go/legacy/exceptions"
	"gitlab.fbs-d.com/dev/go/legacy/i18n"
	"net/http"
)

const (
	RabbitConnectionException       = 12000
	RabbitChannelCreationException  = 12001
	RabbitExchangeCreationException = 12002
	RabbitQueueCreationException    = 12003
	RabbitConsumeException          = 12004
	RabbitPublishException          = 12005
	RabbitQueueBindException        = 12006
	RabbitTriggerEventException     = 12007
	RabbitInvalidParamException     = 12008
	RabbitTimeoutException          = 12009
	RabbitNoConnectionException     = 12010
	RabbitInternalErrorException    = 12011
	RabbitQueueUnbindException      = 12012
	RabbitQueueDeleteException      = 12013
)

func NewRabbitNoConnectionException() excp.IException {
	err := errors.New("Not connected")
	t := i18n.T(i18nCategory, "error.no-connection", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitNoConnectionException,
		http.StatusInternalServerError,
	)
}

func NewRabbitInvalidParamException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.invalid-param", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitInvalidParamException,
		http.StatusInternalServerError,
	)
}

func NewRabbitConnectionException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.connection", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitConnectionException,
		http.StatusInternalServerError,
	)
}

func NewRabbitExchangeCreationException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.create-exchange", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitExchangeCreationException,
		http.StatusInternalServerError,
	)
}

func NewRabbitQueueCreationException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.create-queue", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitQueueCreationException,
		http.StatusInternalServerError,
	)
}

func NewRabbitChannelCreationException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.create-channel", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitChannelCreationException,
		http.StatusInternalServerError,
	)
}

func NewRabbitQueueBindException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.bind-queue", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitQueueBindException,
		http.StatusInternalServerError,
	)
}

func NewRabbitQueueUnbindException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.unbind-queue", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitQueueUnbindException,
		http.StatusInternalServerError,
	)
}

func NewRabbitConsumeException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.consume", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitConsumeException,
		http.StatusInternalServerError,
	)
}

func NewRabbitQueueDeleteException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.delete-queue", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitQueueDeleteException,
		http.StatusInternalServerError,
	)
}
func NewRabbitPublishException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.publish", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitPublishException,
		http.StatusInternalServerError,
	)
}

func NewRabbitTriggerEventException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.trigger-event", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitTriggerEventException,
		http.StatusInternalServerError,
	)
}

func NewRabbitTimeoutException() excp.IException {
	err := errors.New("Request timeout")
	t := i18n.T(i18nCategory, "error.request-timeout", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitTimeoutException,
		http.StatusGatewayTimeout,
	)
}

func NewRabbitInternalException(err error) excp.IException {
	t := i18n.T(i18nCategory, "error.internal", "")
	return excp.NewException(
		excp.NewTranslatedError(err, t, t),
		RabbitInternalErrorException,
		http.StatusGatewayTimeout,
	)
}
