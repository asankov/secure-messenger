package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/asankov/secure-messenger/internal/secretstore"
)

func getKey(stdErr io.Writer) (string, error) {
	if secretKey != "" {
		write(stdErr, "WARN: Passing secret key as a CLI argument is insecure, because the key will be visible in the shell history. Consider using the keychain instead.")

		return secretKey, nil
	}

	if secretKeyFile != "" {
		write(stdErr, "WARN: Storing the secret key in a plain-text file is insecure, because the key can be read by anyone. Consider using the keychain instead.")

		file, err := os.Open(secretKeyFile)
		if err != nil {
			return "", fmt.Errorf("error while opening secret-key file [%s]: %w", secretKeyFile, err)
		}
		secretKeyBytes, err := io.ReadAll(file)
		if err != nil {
			return "", fmt.Errorf("error while reading from secret-key file [%s]: %w", secretKeyFile, err)
		}
		return strings.TrimSpace(string(secretKeyBytes)), nil
	}

	store, err := secretstore.NewKeychainStore()
	if err != nil {
		return "", fmt.Errorf("error while creating keychain store: %w", err)
	}

	secretKey, err := store.GetSecretKey()
	if err != nil {
		return "", fmt.Errorf("error while retrieving secret key from keychain: %w", err)
	}
	return secretKey, nil
}

func write(out io.Writer, msg string) {
	_, _ = out.Write([]byte(msg + "\n"))
}
