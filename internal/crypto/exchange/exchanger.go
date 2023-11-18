package exchange

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asankov/secure-messenger/internal/crypto"
)

type Exchanger struct {
	remoteAddr string
}

func NewExchanger(addr string) *Exchanger {
	return &Exchanger{remoteAddr: addr}
}

type publicKeyExchange struct {
	PublicKey []byte `json:"publicKey"`
}

type secretKeyExchange struct {
	EncryptedSecretKey string `json:"encryptedSecretKey"`
}

func (e *Exchanger) ExchangeSecretKey(privateKey *ecdsa.PrivateKey, secretKey string) error {
	publicKey := privateKey.PublicKey

	marshalledPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return fmt.Errorf("error while marshalling public key: %w", err)
	}

	request, err := json.Marshal(publicKeyExchange{PublicKey: marshalledPublicKey})
	if err != nil {
		return fmt.Errorf("error while marshalling request: %w", err)
	}

	url := e.remoteAddr + "/pub-key-exchange"
	resp, err := http.Post(url, "application/json", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("error while making HTTP call to [%s]: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got non-200 status code from [%s]: %d", url, resp.StatusCode)
	}

	var response publicKeyExchange
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error while decoding response: %w", err)

	}

	publicKeyOtherSide, err := x509.ParsePKIXPublicKey(response.PublicKey)
	if err != nil {
		return fmt.Errorf("error while decoding public key: %w", err)
	}
	publicKeyOtherSideECDSA, ok := publicKeyOtherSide.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("received public key is not of type ecdsa.PublicKey")
	}

	x, _ := publicKeyOtherSideECDSA.Curve.ScalarMult(publicKeyOtherSideECDSA.X, publicKeyOtherSideECDSA.Y, privateKey.D.Bytes())

	sharedSecret := sha256.Sum256(x.Bytes())

	// encrypt the value with the shared key
	encryptor, err := crypto.NewEncryptor(string(sharedSecret[:]))
	if err != nil {
		return fmt.Errorf("error while creating encryptor: %w", err)
	}
	encryptedSecretKey, err := encryptor.Encrypt(secretKey)
	if err != nil {
		return fmt.Errorf("error while encrypting secret key: %w", err)
	}

	request, err = json.Marshal(secretKeyExchange{EncryptedSecretKey: encryptedSecretKey})
	if err != nil {
		return fmt.Errorf("error while marshalling request: %w", err)
	}

	url = e.remoteAddr + "/secret-key-exchange"
	resp, err = http.Post(url, "application/json", bytes.NewReader(request))
	if err != nil {
		return fmt.Errorf("error while making HTTP call to [%s]: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got non-200 status code from [%s]: %d", url, resp.StatusCode)
	}

	return nil
}
