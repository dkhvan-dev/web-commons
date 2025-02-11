package config

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
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

		standardErrMsgsDir := filepath.Join(GetErrMsgPath(), "../errors/msgs")
		files, err := os.ReadDir(standardErrMsgsDir)
		if err != nil {
			err = fmt.Errorf("failed to read localization directory: %w", err)
			return
		}

		for _, file := range files {
			if _, loadErr := Bundle.LoadMessageFile(filepath.Join(standardErrMsgsDir, file.Name())); loadErr != nil {
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

func GetErrMsgPath() string {
	_, fileName, _, _ := runtime.Caller(1)
	return filepath.Dir(fileName)
}
