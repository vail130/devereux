package devereux

import (
	"bufio"
	"fmt"
	"os"
)

func promptForKeyIfNecessary(key string) string {
	var err error
	for key == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter repository key> ")
		key, err = reader.ReadString('\n')
		if err != nil {
			key = ""
			continue
		}
	}

	return key
}

func GetPassword(passwordName string, key string, repoName string) (string, error) {
	key = promptForKeyIfNecessary(key)

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

func SetPassword(passwordName string, password string, key string, repoName string) error {
	key = promptForKeyIfNecessary(key)

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

func CreateRepository(name string, key string, setAsDefault bool) (string, error) {
	key = promptForKeyIfNecessary(key)

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
