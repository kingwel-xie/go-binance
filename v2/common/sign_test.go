package common

import (
	"crypto"
	"crypto/ed25519"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHmac(t *testing.T) {
	assert := assert.New(t)
	secret, data := "secret-key", "timestamp=1234567890&symbol=BNBUSDT"

	got, err := Hmac(secret, data)
	assert.NoError(err)
	assert.NotNil(got)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(data))
	want := fmt.Sprintf("%x", mac.Sum(nil))
	assert.Equal(want, *got)
}

func TestSignFunc(t *testing.T) {
	assert := assert.New(t)

	for _, kt := range []string{KeyTypeHmac, KeyTypeRsa, KeyTypeEd25519} {
		sf, err := SignFunc(kt)
		assert.NoError(err, "SignFunc(%s)", kt)
		assert.NotNil(sf)
	}

	_, err := SignFunc("unknown")
	assert.Error(err)
}

func TestRsa(t *testing.T) {
	assert := assert.New(t)
	const data = "timestamp=1234567890&symbol=BNBUSDT"

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(err)

	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	assert.NoError(err)
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	signed, err := Rsa(string(pemKey), data)
	assert.NoError(err)
	assert.NotNil(signed)

	hashed := sha256.Sum256([]byte(data))
	sig, err := base64.StdEncoding.DecodeString(*signed)
	assert.NoError(err)

	err = rsa.VerifyPKCS1v15(&privateKey.PublicKey, crypto.SHA256, hashed[:], sig)
	assert.NoError(err)
}

func TestEd25519(t *testing.T) {
	assert := assert.New(t)
	const data = "timestamp=1234567890&symbol=BNBUSDT"

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	assert.NoError(err)

	der, err := x509.MarshalPKCS8PrivateKey(priv)
	assert.NoError(err)
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})

	signed, err := Ed25519(string(pemKey), data)
	assert.NoError(err)
	assert.NotNil(signed)

	sig, err := base64.StdEncoding.DecodeString(*signed)
	assert.NoError(err)
	assert.True(ed25519.Verify(pub, []byte(data), sig))
}

func TestRsaInvalidKey(t *testing.T) {
	assert := assert.New(t)
	_, err := Rsa("not-a-pem", "data")
	assert.Error(err)
}

func TestEd25519InvalidKey(t *testing.T) {
	assert := assert.New(t)
	_, err := Ed25519("not-a-pem", "data")
	assert.Error(err)
}
