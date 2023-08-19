package signer_test

import (
	"crypto/rsa"
	"testing"

	"mender/pkg/signer"

	"github.com/stretchr/testify/suite"
)

type SignerTestSuite struct {
	suite.Suite
}

func TestSignerSuite(t *testing.T) {
	suite.Run(t, new(SignerTestSuite))
}

const (
	PrivPem = "-----BEGIN PRIVATE KEY-----\n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\n4aDn+w97eJDyF2Q9uSML/Bj3TqaORw==\n-----END PRIVATE KEY-----\n"
	PublicPem = "-----BEGIN PUBLIC KEY-----\nMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAuTZtoOOsVymtVLKB7zoB\nX5satakTOtmchujCcjPMtRSZqdh9PeM+MK1XyQk8STZFnROnc/3f0Xz6s4zE3fZN\nQdzCgG9cz7FU6V9duFcIaXnypyfOvSApULXEeiZBiOkP5K1gLBul3leGsttpfC2p\niHk8rTyOARucnWJVnj648fauVHFrdo/YM3c+c2crQ883YZELZfHsXrmdOe/63KQD\ndhWeOjLZFhSL9NGzMHs5DknwMX0MoYJqDstexa1K9n33ppF3zlTyp7Y3APN+tDYk\nwqd5Bz9SWT5XT0+qagNdSh/fkrfv+JwurZdRAN5QkVTPcWACBp9sBKabexs9MI3x\nxW9Z7/6dU1tZqiqri6kpjMAjLH7sR49fPhqw12hBbZVgAfeXWhfd1tmCIUL+aK9J\nhkq8BNRHMzV6lr5RqOAba9VWuh2Iz5EyG29tkV8SjNruMnHFJGqwFSU9hBzlSikX\nUOVpON9jITW5qe0SqhorL/zx8rUOsDOrUn8pxorCBJ7DSqjIB5x2ThkAZtgbP0Zj\n5ZzpDDRdQKgz3pXh4dvoLnvNm7uGgPswcr8qYY2XBM30f1QRWgXcmgvNgxzdRXRE\nmSwQMVya5YDEDoo+9IGPgscJL6n28jlFTxtPsku6f4WlRWx66J9n+iKwXr+QvR8R\n/vkwl+UY464mFNM61JJN1X0CAwEAAQ==\n-----END PUBLIC KEY-----\n"
)

func (ks *SignerTestSuite) TestNew() {
	t := ks.Require()

	// No key
	{
		store, err := signer.New()
		t.Nil(store)
		t.ErrorIs(err, signer.ErrNoKey)
	}

	// Key from file
	{
		store, err := signer.New(signer.WithPrivKey([]byte(PrivPem)))
		t.Nil(err)
		t.NotNil(store)
	}
}

func (ks *SignerTestSuite) TestSign() {
	t := ks.Require()
	// Prepare store
	store, err := signer.New(signer.WithPrivKey([]byte(PrivPem)))
	// Verify ctor
	t.Nil(err)
	t.NotNil(store)

	data := []byte("fooobar")
	// Correct sign
	sig, err := store.Sign(data)
	t.Nil(err)

	t.Nil(store.Verify(data, sig))

	// Wrong data
	wrongData := []byte("fooobar2")
	err = store.Verify(wrongData, sig)
	t.ErrorIs(err, rsa.ErrVerification)
	t.NotNil(err)

	// Wrong sig
	// random value in random place
	sig[13] = 5
	err = store.Verify(data, sig)
	t.ErrorIs(err, rsa.ErrVerification)
}
