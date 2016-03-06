package devereux

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

var DEBUG = true

func Debugf(format string, args ...interface{}) {
	if DEBUG {
		log.Printf("DEBUG " + format, args...)
	}
}

func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

///////////////////////////////////////////

var basePath = path.Join(os.Getenv("HOME"), ".devereux/")
var repoPath = path.Join(basePath, "repos/")
var CONFIG_PATH = path.Join(basePath, "config.yaml")

type Config struct {
	CurrentRepository string `yaml:"current_repository"`
}

func (c *Config) Save(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(configPath, os.O_CREATE | os.O_WRONLY, 0600)
	if err != nil {
		Debugf("Problem updating config file")
		return err
	}
	defer f.Close()

	f.Write(data)
	return nil
}

func loadConfig(configPath string) (*Config, error) {
	if exists, _ := fileExists(configPath); !exists {
		return &Config{}, nil
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config *Config
	yaml.Unmarshal(configData, &config)
	return config, nil
}

type PasswordEntry struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type Repository struct {
	Name    string `yaml:"name"`
	Entries []PasswordEntry `yaml:"entries"`

	config  *Config
}

func (r *Repository) Create(setAsDefault bool, configPath string) error {
	filePath := path.Join(repoPath, r.Name)

	if exists, _ := fileExists(filePath); exists {
		return errors.New("Password repository already exists.")
	}

	f, err := os.OpenFile(filePath, os.O_CREATE, 0600)
	if err == nil {
		f.Close()
	}

	if setAsDefault || r.config.CurrentRepository == "" {
		r.config.CurrentRepository = r.Name
		r.config.Save(configPath)
	}

	return err
}

func (r *Repository) GetPassword(passwordName string) (string, error) {
	filePath := path.Join(repoPath, r.Name)

	if exists, _ := fileExists(filePath); !exists {
		return "", errors.New("Password repository does not exist.")
	}

	repoData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	if len(repoData) > 0 {
		yaml.Unmarshal(repoData, &r)
	}

	for i := 0; i < len(r.Entries); i++ {
		entry := r.Entries[i]
		if entry.Name == passwordName {
			return entry.Password, nil
		}
	}

	return "", errors.New("Password name is invalid.")
}

func (r *Repository) SetPassword(passwordName string, password string) error {
	filePath := path.Join(repoPath, r.Name)

	if exists, _ := fileExists(filePath); !exists {
		return errors.New("Password repository does not exist.")
	}

	repoData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if len(repoData) > 0 {
		yaml.Unmarshal(repoData, &r)
	}

	nameUpdated := false
	for i := 0; i < len(r.Entries); i++ {
		entry := r.Entries[i]
		if entry.Name == passwordName {
			r.Entries[i] = PasswordEntry{
				Name: passwordName,
				Password: password,
			}
			nameUpdated = true
		}
	}

	if !nameUpdated {
		r.Entries = append(r.Entries, PasswordEntry{
			Name: passwordName,
			Password: password,
		})
	}

	repoData, err = yaml.Marshal(&r)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, repoData, 0600)
	if err != nil {
		return err
	}

	return nil
}

func GetPassword(passwordName string, repoName string) (string, error) {
	os.MkdirAll(repoPath, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return "", err
	}

	if repoName == "" {
		repoName = config.CurrentRepository
	}

	repo := &Repository{Name: repoName, config: config}
	return repo.GetPassword(passwordName)
}

func SetPassword(passwordName string, password string, repoName string) error {
	os.MkdirAll(repoPath, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return err
	}

	if repoName == "" {
		repoName = config.CurrentRepository
	}

	repo := &Repository{Name: repoName, config: config}
	return repo.SetPassword(passwordName, password)
}

func CreateRepository(name string, setAsDefault bool) (*Repository, error) {
	os.MkdirAll(repoPath, 0777)
	config, err := loadConfig(CONFIG_PATH)
	if err != nil {
		return nil, err
	}

	repo := &Repository{Name: name, config: config}
	err = repo.Create(setAsDefault, CONFIG_PATH)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
