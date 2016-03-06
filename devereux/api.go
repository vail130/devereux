package devereux

import (
	"errors"
	"io/ioutil"
	"os"
)

func promptForKeyIfNecessary(key string, prompt string, envVariable string) (string, error) {
	if key != "" {
		return key, nil
	}

	key = os.Getenv(envVariable)
	if key != "" {
		return key, nil
	}

	var err error
	for key == "" {
		key, err = promptUserForInput(prompt)
		if err != nil {
			return "", err
		}
	}

	return key, nil
}

func GetPassword(passwordName string, repoName string, key string) (string, error) {
	key, err := promptForKeyIfNecessary(key, "Enter repository key> ", "DVRX_KEY")
	if err != nil {
		return "", err
	}

	os.MkdirAll(REPO_PATH, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return "", err
	}

	if repoName == "" {
		repoName = config.CurrentRepository
	}

	if repoName == "" {
		return "", errors.New("No default repository found. Specify one with the -r flag.")
	}

	repo := &Repository{Name: repoName, config: config}
	return repo.GetPassword(passwordName, key)
}

func SetPassword(passwordName string, repoName string, key string, password string) error {
	key, err := promptForKeyIfNecessary(key, "Enter repository key> ", "DVRX_KEY")
	if err != nil {
		return err
	}

	password, err = promptForKeyIfNecessary(password, "Enter password> ", "DVRX_PASS")
	if err != nil {
		return err
	}

	os.MkdirAll(REPO_PATH, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return err
	}

	if repoName == "" {
		repoName = config.CurrentRepository
	}

	if repoName == "" {
		return "", errors.New("No default repository found. Specify one with the -r flag.")
	}

	repo := &Repository{Name: repoName, config: config}
	return repo.SetPassword(passwordName, password, key)
}

func CreateRepository(name string, setAsDefault bool, key string) (string, error) {
	key, err := promptForKeyIfNecessary(key, "Enter repository key> ", "DVRX_KEY")
	if err != nil {
		return "", err
	}

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

func GetRepositories() ([]string, error) {
	names := make([]string, 0)

	files, err := ioutil.ReadDir(REPO_PATH)
	if err != nil {
		return names, err
	}

	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

func DeleteRepository(name string) error {
	repo := &Repository{Name: name}
	err := repo.Delete()
	if err != nil {
		return err
	}

	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return err
	}

	if config.CurrentRepository == name {
		config.CurrentRepository = ""
		config.Save(CONFIG_PATH)
	}

	return nil
}
