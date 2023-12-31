package secretstore

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	service     = "secure-messenger"
	accessGroup = "org.secure-messenger"

	location = "keychain"

	secretKeyKey = "secret-key"
)

type KeychainStore struct{}

func NewKeychainStore() (*KeychainStore, error) {
	return &KeychainStore{}, nil
}

func (s *KeychainStore) StoreSecretKey(key string) (string, error) {
	return s.Store(secretKeyKey, key)
}

func (s *KeychainStore) GetSecretKey() (string, error) {
	return s.Get(secretKeyKey)
}

func (s *KeychainStore) Store(key, value string) (string, error) {
	err := keyring.Set(service, key, value)
	if err != nil {
		return "", fmt.Errorf("error while storing [%s] in keyring: %w", key, err)
	}

	return location, nil
}

func (s *KeychainStore) Get(key string) (string, error) {
	value, err := keyring.Get(service, key)
	if err != nil {
		return "", fmt.Errorf("error while looking up [%s] in keyring: %w", key, err)
	}
	return value, nil

}
