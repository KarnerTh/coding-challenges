package crypto

import (
	"fmt"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/errors"
)

type SignatureAlgorithm string

const (
	ECC SignatureAlgorithm = "ECC"
	RSA SignatureAlgorithm = "RSA"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
	Verify(dataToBeVerified []byte, signature []byte) bool
}

func CreateSigner(algorithm SignatureAlgorithm) (Signer, error) {
	switch algorithm {
	case RSA:
		return NewRsaSigner()
	case ECC:
		return NewEccSigner()
	default:
		return nil, errors.BadInputError{Msg: fmt.Sprintf("algorithm not supporetd: %s", algorithm)}
	}
}
