package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

type Signer struct {
	priv *rsa.PrivateKey
}

type Option func(k *Signer) error

var (
	ErrNoKey = errors.New("no key provided")
)

func New(options ...Option) (*Signer, error) {
	k := &Signer{}
	for _, opt := range options {
		if err := opt(k); err != nil {
			return nil, err
		}
	}
	if k.priv == nil {
		return nil, ErrNoKey
	}

	return k, nil
}

func (k *Signer) PublicKeyPEM() string {
	pubBytes, err := x509.MarshalPKIXPublicKey(&k.priv.PublicKey)
	if err != nil {
		panic(err)
	}
	pubPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: pubBytes,
		},
	)

	return string(pubPem)
}

func (k *Signer) Sign(data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, k.priv, crypto.SHA256, hashed[:])
}

func (k *Signer) Verify(data []byte, sig []byte) error {
	cr := crypto.SHA256.New()
	cr.Write(data)
	return rsa.VerifyPKCS1v15(&k.priv.PublicKey, crypto.SHA256, cr.Sum(nil), sig)
}

func PEMToPrivate(privPem []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPem)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPriv, ok := priv.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not a rsa.PrivateKey")
	}

	return rsaPriv, nil
}

func PEMToPublic(pubPem []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPem)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not a rsa.PublicKey")
	}
	return rsaKey, nil
}
