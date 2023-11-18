package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func getKey() (string, error) {
	if secretKey != "" {
		return secretKey, nil
	}

	if secretKeyFile != "" {
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

	return "", fmt.Errorf(`either "secret-key" or "secret-key-file" must be specified`)
}
