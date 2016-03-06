package devereux

import (
	"fmt"
	"os"

	"github.com/howeyc/gopass"
)

func promptForKeyIfNecessary(key string, prompt string, envVariable string) string {
	key = os.Getenv(envVariable)
	if key != "" {
		return key
	}

	var err error
	var keyBytes []byte

	for key == "" {
		fmt.Printf("Enter %s> ", prompt)
		keyBytes, err = gopass.GetPasswd()
		if err == nil {
			key = string(keyBytes)
		}
	}

	return key
}

func GetPassword(passwordName string, repoName string, key string) (string, error) {
	key = promptForKeyIfNecessary(key, "repository key", "DVRX_KEY")

	os.MkdirAll(REPO_PATH, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return "", err
	}

	if repoName == "" {
		repoName = config.CurrentRepository
	}

	repo := &Repository{Name: repoName, config: config}
	return repo.GetPassword(passwordName, key)
}

func SetPassword(passwordName string, repoName string, key string, password string) error {
	key = promptForKeyIfNecessary(key, "repository key", "DVRX_KEY")
	password = promptForKeyIfNecessary(password, "password", "DVRX_PASS")

	os.MkdirAll(REPO_PATH, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return err
	}

	if repoName == "" {
		repoName = config.CurrentRepository
	}

	repo := &Repository{Name: repoName, config: config}
	return repo.SetPassword(passwordName, password, key)
}

func CreateRepository(name string, setAsDefault bool, key string) (string, error) {
	key = promptForKeyIfNecessary(key, "repository key", "DVRX_KEY")

	os.MkdirAll(REPO_PATH, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return "", err
	}

	repo := &Repository{Name: name, config: config}
	err = repo.Create(key, setAsDefault, CONFIG_PATH)
	if err != nil {
		return "", err
	}

	return repo.Name, nil
}
