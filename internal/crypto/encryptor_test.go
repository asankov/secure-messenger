package crypto_test

import (
	"testing"

	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/stretchr/testify/require"
)

const (
	msg = "A purely peer-to-peer version of electronic cash would allow online payments to be sent directly from one party to another without going through a financial institution."
)

// TestEncryptor tests that a message encrypted with one Encryptor (hence, one secret key)
// can be decrypted back to its original form with the same Encryptor (hence, same secret key).
// Since we are supporting 16, 24 and 32 byte keys, the test runs one for a key with each size.
func TestEncryptor(t *testing.T) {
	testCases := []struct {
		name string
		key  string
	}{
		{
			name: "16-byte key",
			key:  "some16bitkeyabcd",
		},
		{
			name: "24-byte key",
			key:  "some24bitkeyabcdefghijkl",
		},
		{
			name: "32-byte key",
			key:  "some32bitkeyabcdefghijklmnopqrst",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			encryptor, err := crypto.NewEncryptor(testCase.key)
			require.NoError(t, err)

			ciphertext, err := encryptor.Encrypt(msg)
			require.NoError(t, err)

			decrypted, err := encryptor.Decrypt(ciphertext)
			require.NoError(t, err)

			require.Equal(t, msg, decrypted)
			// test that the unit tests in GH actions will fail
			require.Equal(t, 0, 1)
		})
	}

}

// TestEncryptorWrongSizeKey verifies that we cannot create an Encryptor with a key size different that 16, 24 or 32.
func TestEncryptorWrongSizeKey(t *testing.T) {
	_, err := crypto.NewEncryptor("abc123")
	require.Error(t, err)
}

// TestEncryptorWrongKey verifies that a message encrypted with one key cannot be decrypted by another key.
func TestEncryptorWrongKey(t *testing.T) {
	secretKey := "some32bitkeyabcdefghijklmnopqrst"
	encryptor, err := crypto.NewEncryptor(secretKey)
	require.NoError(t, err)

	ciphertext, err := encryptor.Encrypt(msg)
	require.NoError(t, err)

	anotherEncryptor, err := crypto.NewEncryptor("another32bitkeyabcdefghijklmnopq")
	require.NoError(t, err)

	_, err = anotherEncryptor.Decrypt(ciphertext)
	require.Error(t, err)
}
