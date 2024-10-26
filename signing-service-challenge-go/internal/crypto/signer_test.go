package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	t.Parallel()
	availableSigners := []struct {
		algoritym SignatureAlgorithm
		signer    func() (Signer, error)
	}{
		{algoritym: RSA, signer: NewRsaSigner},
		{algoritym: ECC, signer: NewEccSigner},
	}

	for _, test := range availableSigners {
		t.Run(string(test.algoritym), func(t *testing.T) {
			t.Parallel()

			t.Run("Valid signature", func(t *testing.T) {
				t.Parallel()
				dataToBeSigned := []byte("data")
				signer, err := test.signer()
				assert.Nil(t, err)

				signature, err := signer.Sign(dataToBeSigned)
				assert.Nil(t, err)
				assert.NotEmpty(t, signature)

				validSignature := signer.Verify(dataToBeSigned, signature)
				assert.True(t, validSignature)
			})

			t.Run("Invalid signature", func(t *testing.T) {
				t.Parallel()
				dataToBeSigned := []byte("data")
				signer, err := test.signer()
				assert.Nil(t, err)

				signature, err := signer.Sign(dataToBeSigned)
				assert.Nil(t, err)
				assert.NotEmpty(t, signature)

				invalidData := []byte("signatureShouldFail")
				validSignature := signer.Verify(invalidData, signature)
				assert.False(t, validSignature)
			})

		})
	}
}
