package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

// GenerateSecretKey generates a secret key of the given size.
//
// Supported sizes are 16, 24 and 32.
func GenerateSecretKey(size int) (string, error) {
	if size != 16 && size != 24 && size != 32 {
		return "", fmt.Errorf("key size of [%d] is not allowed. Allowed values are 16, 24 and 32.", size)
	}

	size = normalizeSize(size)
	key := make([]byte, size)

	n, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("error while generating secret key: %w", err)
	} else if n != size {
		// according to the documentation of rand.Read we should never get here,
		// because if n != len(b) then err will be non-nil,
		// but let's do this check just to be safe.
		return "", errors.New("unable to generate key with the desired size")
	}

	return base64.StdEncoding.EncodeToString(key), nil
}

// normalizeSize changes the size of the desired key,
// that is needed, because the base64 encoding we do at the end
// adds additional bytes to the key.
func normalizeSize(size int) int {
	if size == 16 {
		return 12
	}
	if size == 24 {
		return 16
	}
	if size == 32 {
		return 24
	}
	// should never get here
	return -1
}
