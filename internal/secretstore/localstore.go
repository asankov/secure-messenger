package secretstore

import (
	"fmt"
	"os"
	"path/filepath"
)

type LocalStore struct{}

func NewLocalStore() (*LocalStore, error) {
	newpath := filepath.Join("~", ".secure-messenger")
	if err := os.MkdirAll(newpath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error while creating directory: %w", err)
	}

	return &LocalStore{}, nil
}

func (s *LocalStore) Store(key, value string) error {
	return os.WriteFile("~/.secure-messenger/secret-"+key, []byte(value), 0600)
}

func (s *LocalStore) Get(key string) (string, error) {
	secret, err := os.ReadFile("~/.secure-messenger/secret-" + key)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}
