package domain

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
)

type SignatureDevice struct {
	Id               string                    `json:"id"`
	Algorithm        crypto.SignatureAlgorithm `json:"algorithm"`
	Label            string                    `json:"label"`
	SignatureCounter int                       `json:"-"`
	LastSignature    string                    `json:"-"`
	mu               sync.Mutex
	signer           crypto.Signer
}

type Signature struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
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
	UpdateSigningMetaInfo(id string, counter int, lastSignature string) (*SignatureDevice, error)
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
		LastSignature: lastSignatureFallback,
		signer:        signer,
	}

	return s.repo.Create(&device)
}

func (s SignatureDeviceService) GetAll() ([]*SignatureDevice, error) {
	return s.repo.GetAll()
}

func (s SignatureDeviceService) GetById(id string) (*SignatureDevice, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("id must be specified")
	}

	return s.repo.GetById(id)
}

func (s SignatureDeviceService) Sign(deviceId string, data string) (*Signature, error) {
	device, err := s.GetById(deviceId)
	if err != nil {
		return nil, err
	}

	device.mu.Lock()
	defer device.mu.Unlock()

	signData := fmt.Sprintf("%d_%s_%s", device.SignatureCounter, data, device.LastSignature)
	slog.Debug("data is being signed", "value", signData)

	signature, err := device.signer.Sign([]byte(signData))
	signatureBase64 := base64.StdEncoding.EncodeToString(signature)

	s.repo.UpdateSigningMetaInfo(deviceId, device.SignatureCounter+1, signatureBase64)

	return &Signature{
		Signature:  signatureBase64,
		SignedData: signData,
	}, nil
}
