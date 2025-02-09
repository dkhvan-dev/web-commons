package errors

import (
	"encoding/json"
	"fmt"
	"github.com/dkhvan-dev/web-commons/constants"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
)

var ErrLocalizer *i18n.Localizer

type CustomError struct {
	Msg  string
	Code int
}

func (c CustomError) Error() string {
	return c.Msg
}

func InitLocalization(lang string) error {
	langTag, err := language.Parse(lang)
	if err != nil {
		panic(err)
	}

	bundle := i18n.NewBundle(langTag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if _, err := bundle.LoadMessageFile("./msgs/errors.en.json"); err != nil {
		return fmt.Errorf("failed to load English translations: %w", err)
	}

	if _, err := bundle.LoadMessageFile("./msgs/errors.ru.json"); err != nil {
		return fmt.Errorf("failed to load Russian translations: %w", err)
	}

	if _, err := bundle.LoadMessageFile("./msgs/errors.kk.json"); err != nil {
		return fmt.Errorf("failed to load Kazakh translations: %w", err)
	}

	ErrLocalizer = i18n.NewLocalizer(bundle, lang)
	return nil
}

func NewLocalizedError(key string, code int) *CustomError {
	message := ErrLocalizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	return &CustomError{message, code}
}

func HandleLocalizedError(w http.ResponseWriter, err *CustomError) {
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
