package domain

import (
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
)

type SignatureDevice struct {
	Id               string                    `json:"id"`
	Algorithm        crypto.SignatureAlgorithm `json:"algorithm"`
	Label            string                    `json:"label"`
	mu               sync.Mutex
	signer           crypto.Signer
	signatureCounter int
	lastSignature    string
}

type SignatureDeviceService struct {
	repo SignatureDeviceRepository
}

func NewSignatureDeviceService(repo SignatureDeviceRepository) SignatureDeviceService {
	return SignatureDeviceService{repo: repo}
}

type SignatureDeviceRepository interface {
	Create(device *SignatureDevice) (*SignatureDevice, error)
	GetAll() ([]*SignatureDevice, error)
	GetById(id string) (*SignatureDevice, error)
}

func (s SignatureDeviceService) Create(id string, algorithm crypto.SignatureAlgorithm, label string) (*SignatureDevice, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("id must be specified")
	}

	signer, err := crypto.CreateSigner(algorithm)
	if err != nil {
		return nil, err
	}

	lastSignatureFallback := base64.StdEncoding.EncodeToString([]byte(id))
	device := SignatureDevice{
		Id:            id,
		Algorithm:     algorithm,
		Label:         label,
		lastSignature: lastSignatureFallback,
		signer:        signer,
	}

	return s.repo.Create(&device)
}

func (r SignatureDeviceService) GetAll() ([]*SignatureDevice, error) {
	return r.repo.GetAll()
}

func (r SignatureDeviceService) GetById(id string) (*SignatureDevice, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("id must be specified")
	}

	return r.repo.GetById(id)
}
