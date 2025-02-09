package config

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"sync"
)

var (
	Bundle *i18n.Bundle
	once   sync.Once
	err    error
)

func InitLocalization(paths []string) error {
	once.Do(func() {
		Bundle = i18n.NewBundle(language.English)
		Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		standardPaths := []string{
			"./errors/msgs/errors.en.json",
			"./errors/msgs/errors.ru.json",
			"./errors/msgs/errors.kk.json",
		}
		for _, path := range standardPaths {
			if _, loadErr := Bundle.LoadMessageFile(path); loadErr != nil {
				err = fmt.Errorf("failed to load standard error messages: %w", err)
				return
			}
		}

		for _, path := range paths {
			if _, loadErr := Bundle.LoadMessageFile(path); loadErr != nil {
				err = fmt.Errorf("failed to load translations '%s': %w", path, loadErr)
				return
			}
		}
	})

	return err
}
