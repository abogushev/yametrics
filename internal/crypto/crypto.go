package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func Encrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		message,
		nil)

	if err != nil {
		return nil, err
	}
	return encryptedBytes, nil
}

func Decrypt(privateKey *rsa.PrivateKey, encryptedMessage []byte) ([]byte, error) {
	return privateKey.Decrypt(nil, encryptedMessage, &rsa.OAEPOptions{Hash: crypto.SHA256})
}

func readPem(fileName string) ([]byte, error) {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(file)
	if block == nil {
		return nil, errors.New("failed parse pem file")
	}
	return block.Bytes, err
}
func ReadPublicKey(fileName string) (*rsa.PublicKey, error) {
	file, err := readPem(fileName)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PublicKey(file)
}

func ReadPrivateKey(fileName string) (*rsa.PrivateKey, error) {
	file, err := readPem(fileName)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(file)
}
