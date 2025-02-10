package config

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"path/filepath"
	"runtime"
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

		_, filename, _, _ := runtime.Caller(0)
		webCommonsPath := filepath.Join(filepath.Dir(filename), "../errors/msgs")
		standardPaths := []string{
			filepath.Join(webCommonsPath, "errors.en.json"),
			filepath.Join(webCommonsPath, "errors.ru.json"),
			filepath.Join(webCommonsPath, "errors.kk.json"),
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
