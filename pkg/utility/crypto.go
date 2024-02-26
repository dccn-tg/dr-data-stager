package utility

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

// EncryptStringWithRsaKey encrypts the `plaintext` with the asymetric key, and
// returns the encrypted data represented in the base64 standard encoding.
func EncryptStringWithRsaKey(plaintext string, keyFile string) (*string, error) {

	publicKeyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)

	if err != nil {
		return nil, err
	}

	encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), []byte(plaintext))
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(encrypted)

	return &encoded, nil

}

// DecryptStringWithRsaKey decrypts the encrypted string encoded in base64 standard encoding,
// and returns the decrypted string.
func DecryptStringWithRsaKey(encryptedBase64 string, keyFile string) (*string, error) {

	privateKeyPEM, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	encrypted, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, err
	}

	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), encrypted)
	if err != nil {
		return nil, err
	}

	s := string(plaintext)

	return &s, nil
}
