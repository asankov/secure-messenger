package exchange

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/asankov/secure-messenger/internal/crypto"
	"github.com/asankov/secure-messenger/internal/secretstore"
)

const (
	sharedSecretStoreKey = "shared-secret"
	secretKeyStoreKey    = "secret-key"
)

type Listener struct {
	logger *slog.Logger
	store  *secretstore.KeychainStore
}

func NewListener() (*Listener, error) {
	store, err := secretstore.NewKeychainStore()
	if err != nil {
		return nil, err
	}
	return &Listener{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
		store:  store,
	}, nil
}

func (l *Listener) SetLogger(logger *slog.Logger) {
	l.logger = logger
}

func (l *Listener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/pub-key-exchange" {
		l.handlePublicKeyExchange(w, r)
		return
	}

	if r.URL.Path == "/secret-key-exchange" {
		l.handleSecretKeyExchange(w, r)
		return
	}

	l.logger.Error("received request on unknown path", "path", r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
}

func (l *Listener) handlePublicKeyExchange(w http.ResponseWriter, r *http.Request) {
	var request publicKeyExchange
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		l.logger.Error("error while generating decoding request", "error", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		l.logger.Error("error while generating private key", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	publicKeyOtherSide, err := x509.ParsePKIXPublicKey(request.PublicKey)
	if err != nil {
		l.logger.Error("error while parsing public key from request", "error", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	publicKeyOtherSideECDSA, ok := publicKeyOtherSide.(*ecdsa.PublicKey)
	if !ok {
		l.logger.Error("received public key is not ecdsa.PublicKey", "error", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	x, _ := publicKeyOtherSideECDSA.Curve.ScalarMult(publicKeyOtherSideECDSA.X, publicKeyOtherSideECDSA.Y, privateKey.D.Bytes())
	sharedSecret := sha256.Sum256(x.Bytes())

	if _, err := l.store.Store(sharedSecretStoreKey, string(sharedSecret[:])); err != nil {
		l.logger.Error("error while storing shared secret", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	// shared secret is successfully computed an saved.
	// now return our public key to the other side
	// so that they can also compute the shared secret.

	publicKey := privateKey.PublicKey

	marshalledPublicKey, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		l.logger.Error("error while marshalling public key", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	response, err := json.Marshal(publicKeyExchange{PublicKey: marshalledPublicKey})
	if err != nil {
		l.logger.Error("error while marshalling response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
	if _, err := w.Write(response); err != nil {
		l.logger.Error("error while writing response", "error", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}

func (l *Listener) handleSecretKeyExchange(w http.ResponseWriter, r *http.Request) {
	var request secretKeyExchange
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		panic(err)
	}

	sharedSecret, err := l.store.Get(sharedSecretStoreKey)
	if err != nil {
		panic(err)
	}

	encryptor, err := crypto.NewEncryptor(sharedSecret)
	if err != nil {
		panic(err)
	}

	secretKey, err := encryptor.Decrypt(request.EncryptedSecretKey)
	if err != nil {
		panic(err)
	}

	l.logger.Info("Retrieved secret key")

	location, err := l.store.StoreSecretKey(secretKey)
	if err != nil {
		l.logger.Error("error while storing secret key", "error", err)
	} else {
		l.logger.Info("Secret key stored", "location", location)
	}

	w.WriteHeader(http.StatusOK)
}
