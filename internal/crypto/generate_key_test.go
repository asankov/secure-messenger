package crypto_test

import (
	"fmt"
	"testing"

	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/stretchr/testify/require"
)

func TestGenerateSecretKey(t *testing.T) {
	testCases := []struct {
		size int
	}{
		{size: 16},
		{size: 24},
		{size: 32},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("size = %d", testCase.size), func(t *testing.T) {
			secretKey, err := crypto.GenerateSecretKey(testCase.size)
			require.NoError(t, err)
			require.Len(t, secretKey, testCase.size)
			require.NotEmpty(t, secretKey)
		})
	}

	t.Run("SizeNotAllowed", func(t *testing.T) {
		_, err := crypto.GenerateSecretKey(15)
		require.Error(t, err)
	})
}
