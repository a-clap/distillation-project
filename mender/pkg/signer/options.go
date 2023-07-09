package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
)

func WithPrivKey(privPem []byte) Option {
	return func(k *Signer) error {
		priv, err := PEMToPrivate(privPem)
		if err != nil {
			return fmt.Errorf("failed to parse priv key from pem: %w", err)
		}
		k.priv = priv
		return nil
	}
}

func WithNewKey() Option {
	return func(k *Signer) error {
		k.priv, _ = rsa.GenerateKey(rand.Reader, 4096)
		return nil
	}
}
