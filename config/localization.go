package config

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"net/http"
)

var Bundle *i18n.Bundle

func InitLocalization() error {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if _, err := Bundle.LoadMessageFile("./msgs/errors.en.json"); err != nil {
		return fmt.Errorf("failed to load English translations: %w", err)
	}

	if _, err := Bundle.LoadMessageFile("./msgs/errors.ru.json"); err != nil {
		return fmt.Errorf("failed to load Russian translations: %w", err)
	}

	if _, err := Bundle.LoadMessageFile("./msgs/errors.kk.json"); err != nil {
		return fmt.Errorf("failed to load Kazakh translations: %w", err)
	}

	return nil
}

func GetLocalizer(r *http.Request) *i18n.Localizer {
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang == "" {
		acceptLang = "en"
	}

	// Создаем локализатор
	return i18n.NewLocalizer(Bundle, acceptLang)
}
