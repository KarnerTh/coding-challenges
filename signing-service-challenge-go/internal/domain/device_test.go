package domain_test

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence"
	"github.com/stretchr/testify/assert"
)

func createService() domain.SignatureDeviceService {
	return domain.NewSignatureDeviceService(persistence.NewSignatureDeviceInMemoryRepo())
}

func TestDeviceCreation(t *testing.T) {
	t.Parallel()
	t.Run("Valid input", func(t *testing.T) {
		t.Parallel()
		s := createService()
		device, err := s.Create("validId", crypto.RSA, "")

		assert.Nil(t, err)
		assert.NotNil(t, device)
		assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("validId")), device.LastSignature)
		assert.Equal(t, 0, device.SignatureCounter)
	})

	t.Run("Empty id is not allowed", func(t *testing.T) {
		t.Parallel()
		s := createService()
		device, err := s.Create("", crypto.RSA, "")

		assert.NotNil(t, err)
		assert.Nil(t, device)
		assert.IsType(t, errors.BadInputError{}, err)
	})

	t.Run("Unsupported algorythm", func(t *testing.T) {
		t.Parallel()
		s := createService()
		device, err := s.Create("validId", crypto.SignatureAlgorithm("XXX"), "")

		assert.NotNil(t, err)
		assert.Nil(t, device)
		assert.IsType(t, errors.BadInputError{}, err)
	})

	t.Run("Id must be unique", func(t *testing.T) {
		t.Parallel()
		s := createService()
		s.Create("validId", crypto.RSA, "")
		device, err := s.Create("validId", crypto.RSA, "")

		assert.NotNil(t, err)
		assert.Nil(t, device)
		assert.IsType(t, errors.ConflictError{}, err)
	})
}

func TestSignatureCreation(t *testing.T) {
	t.Parallel()
	t.Run("Valid signature", func(t *testing.T) {
		t.Parallel()
		s := createService()
		_, err := s.Create("validId", crypto.RSA, "")
		assert.Nil(t, err)
		signature, err := s.Sign("validId", "dataToBeSigned")

		assert.Nil(t, err)
		assert.NotNil(t, signature)
		assert.Equal(t, "0", strings.Split(signature.SignedData, "_")[0])
	})

	t.Run("Signature counter increments and lastSignature is used", func(t *testing.T) {
		t.Parallel()
		s := createService()
		_, err := s.Create("validId", crypto.RSA, "")
		assert.Nil(t, err)
		signature1, _ := s.Sign("validId", "dataToBeSigned1")
		signature2, _ := s.Sign("validId", "dataToBeSigned2")

		assert.Equal(t, "0", strings.Split(signature1.SignedData, "_")[0])
		assert.Equal(t, "1", strings.Split(signature2.SignedData, "_")[0])
		assert.Equal(t, signature1.Signature, strings.Split(signature2.SignedData, "_")[2])
	})
}
