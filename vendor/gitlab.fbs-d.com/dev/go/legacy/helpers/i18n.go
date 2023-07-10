package toolkit

import "gitlab.fbs-d.com/dev/go/legacy/i18n"

const i18nCategory = "helpers"

var messages = `{
		"error.generate-random-string-error":"Generate random string error",
		"error.nil-data-error":"Can not convert nil data",
		"error.currency-conversion":"Currency conversion error",
		"error.parse-time":"Time parsing error",
		"error.decode-image": "Decoding image error",
		"error.save-image": "Saving image error"
}`

func LoadDefaultTranslations() error {
	return i18n.LoadRaw(messages, i18nCategory, "")
}
