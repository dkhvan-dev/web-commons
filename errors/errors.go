package errors

import (
	"encoding/json"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

type CustomError struct {
	Msg  string
	Code int
}

func (c CustomError) Error() string {
	return c.Msg
}

func NewLocalizedError(localizer *i18n.Localizer, key string, code int) *CustomError {
	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	return &CustomError{message, code}
}

func HandleLocalizedError(w http.ResponseWriter, r *http.Request, key string, code int) {
	localizer := config.GetLocalizer(r)
	err := NewLocalizedError(localizer, key, code)

	errJson := struct {
		Error string `json:"error"`
	}{
		Error: err.Error(),
	}

	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	w.WriteHeader(err.Code)

	data, marshalErr := json.Marshal(errJson)
	if marshalErr != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(data); writeErr != nil {
		http.Error(w, `{"error": "internal server error"}`, http.StatusInternalServerError)
		return
	}
}
