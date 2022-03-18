package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

//DecryptSecret decrypts a given secret, returns the string if only a plaintext secre is given
func DecryptSecret(secret string) (string, error) {
	data := strings.Split(secret, ":")

	if data[0] == "rsa" {
		pemPrivate, ok := os.LookupEnv("RSA_PRIVATE_KEY")
		if !ok {
			return "", fmt.Errorf("environment variable RSA_PRIVATE_KEY is not set, will not be able to decrypt secrets")
		}

		secret, err := decryptPrivatePEM(data[1], pemPrivate)
		if err != nil {
			return "", err
		}
		return secret, nil
	}
	return secret, nil
}

func decryptPrivatePEM(message string, keyfile string) (string, error) {
	key, err := ioutil.ReadFile(keyfile) // just pass the file name
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(key)

	if err != nil {
		return "", err
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}
	ct, err := rsaOaepDecrypt(message, *privateKey)
	if err != nil {
		return "", err
	}

	return ct, nil
}

func rsaOaepDecrypt(cipherText string, privKey rsa.PrivateKey) (string, error) {
	ct, _ := base64.StdEncoding.DecodeString(cipherText)
	label := []byte("OAEP Encrypted")
	rng := rand.Reader

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, ct, label)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
