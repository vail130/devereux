package devereux

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"bytes"
	"crypto/rand"
	"gopkg.in/yaml.v2"
	"math"
)

var BASE_PATH = path.Join(os.Getenv("HOME"), ".devereux/")
var REPO_PATH = path.Join(BASE_PATH, "repos/")
var BOUNDARY_BYTES = []byte("私のホバークラフトは鰻でいっぱいです")
var MIN_REPO_BYTE_LENGTH = 5 * int(math.Pow(10, float64(6)))

type PasswordEntry struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type Repository struct {
	Name    string          `yaml:"name"`
	Entries []PasswordEntry `yaml:"entries"`

	config *Config
}

func (r *Repository) LoadFromFile(key string) error {
	filePath := path.Join(REPO_PATH, r.Name)

	if exists, _ := fileExists(filePath); !exists {
		return errors.New("Password repository does not exist.")
	}

	repoData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(repoData) == 0 {
		return errors.New("Repository is invalid. No data in repository.")
	}

	repoData, err = decrypt(key, repoData)
	if err != nil {
		return err
	}

	boundaryIndex := bytes.LastIndex(repoData, BOUNDARY_BYTES)
	if boundaryIndex == -1 {
		return errors.New("Repository is invalid. Cannot find boundary.")
	}

	repoData = repoData[:boundaryIndex]
	yaml.Unmarshal(repoData, &r)

	return nil
}

func (r *Repository) WriteToFile(key string) error {
	repoData, err := yaml.Marshal(&r)
	if err != nil {
		return err
	}

	if len(repoData) > MIN_REPO_BYTE_LENGTH-len(BOUNDARY_BYTES) {
		return errors.New("Your repository is too big. Allocate more size for it (future feature).")
	}

	numBufferBytes := MIN_REPO_BYTE_LENGTH - len(BOUNDARY_BYTES) - len(repoData)
	bufferByteSlice := make([]byte, numBufferBytes)
	_, err = rand.Read(bufferByteSlice)
	if err != nil {
		return err
	}

	repoData = concatByteSlices(repoData, BOUNDARY_BYTES, bufferByteSlice)

	filePath := path.Join(REPO_PATH, r.Name)
	repoData, err = encrypt(key, repoData)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, repoData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPassword(passwordName string, key string) (string, error) {
	err := r.LoadFromFile(key)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(r.Entries); i++ {
		entry := r.Entries[i]
		if entry.Name == passwordName {
			return entry.Password, nil
		}
	}

	return "", errors.New("Password name is invalid.")
}

func (r *Repository) SetPassword(passwordName string, password string, key string) error {
	err := r.LoadFromFile(key)
	if err != nil {
		return err
	}

	nameUpdated := false
	for i := 0; i < len(r.Entries); i++ {
		entry := r.Entries[i]
		if entry.Name == passwordName {
			r.Entries[i] = PasswordEntry{
				Name:     passwordName,
				Password: password,
			}
			nameUpdated = true
		}
	}

	if !nameUpdated {
		r.Entries = append(r.Entries, PasswordEntry{
			Name:     passwordName,
			Password: password,
		})
	}

	return r.WriteToFile(key)
}

func (r *Repository) Create(key string, setAsDefault bool, configPath string) error {
	filePath := path.Join(REPO_PATH, r.Name)

	if exists, _ := fileExists(filePath); exists {
		return errors.New("Password repository already exists.")
	}

	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	err = r.WriteToFile(key)
	if err != nil {
		return err
	}

	if setAsDefault || r.config.CurrentRepository == "" {
		r.config.CurrentRepository = r.Name
		r.config.Save(configPath)
	}

	return err
}
