package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

type eccSigner struct {
	keyPair ECCKeyPair
}

func NewEccSigner() (Signer, error) {
	keyGen := ECCGenerator{}
	keyPair, err := keyGen.Generate()

	if err != nil {
		return nil, err
	}

	return &eccSigner{
		keyPair: *keyPair,
	}, nil
}

func (s *eccSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hashed := sha256.Sum256(dataToBeSigned)
	return ecdsa.SignASN1(rand.Reader, s.keyPair.Private, hashed[:])
}

func (s *eccSigner) Verify(dataToBeVerified []byte, signature []byte) bool {
	hashed := sha256.Sum256(dataToBeVerified)
	return ecdsa.VerifyASN1(s.keyPair.Public, hashed[:], signature)
}

// ECCKeyPair is a DTO that holds ECC private and public keys.
type ECCKeyPair struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

// ECCMarshaler can encode and decode an ECC key pair.
type ECCMarshaler struct{}

// NewECCMarshaler creates a new ECCMarshaler.
func NewECCMarshaler() ECCMarshaler {
	return ECCMarshaler{}
}

// Encode takes an ECCKeyPair and encodes it to be written on disk.
// It returns the public and the private key as a byte slice.
func (m ECCMarshaler) Encode(keyPair ECCKeyPair) ([]byte, []byte, error) {
	privateKeyBytes, err := x509.MarshalECPrivateKey(keyPair.Private)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(keyPair.Public)
	if err != nil {
		return nil, nil, err
	}

	encodedPrivate := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE_KEY",
		Bytes: privateKeyBytes,
	})

	encodedPublic := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC_KEY",
		Bytes: publicKeyBytes,
	})

	return encodedPublic, encodedPrivate, nil
}

// Decode assembles an ECCKeyPair from an encoded private key.
func (m ECCMarshaler) Decode(privateKeyBytes []byte) (*ECCKeyPair, error) {
	block, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return &ECCKeyPair{
		Private: privateKey,
		Public:  &privateKey.PublicKey,
	}, nil
}