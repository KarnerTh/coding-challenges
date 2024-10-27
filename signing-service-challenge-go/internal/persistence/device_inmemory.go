package persistence

import (
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/errors"
)

type SignatureDeviceInMemoryRepo struct {
	devices sync.Map
}

func NewSignatureDeviceInMemoryRepo() *SignatureDeviceInMemoryRepo {
	return &SignatureDeviceInMemoryRepo{}
}

func (r *SignatureDeviceInMemoryRepo) Create(device *domain.SignatureDevice) (*domain.SignatureDevice, error) {
	_, ok := r.devices.Load(device.Id)
	if ok {
		return nil, errors.ConflictError{Id: device.Id}
	}

	r.devices.Store(device.Id, device)
	return device, nil
}

func (r *SignatureDeviceInMemoryRepo) GetAll() ([]*domain.SignatureDevice, error) {
	var devices []*domain.SignatureDevice
	var err error

	r.devices.Range(func(key, value any) bool {
		device, ok := value.(*domain.SignatureDevice)
		if !ok {
			err = errors.InternalError{Msg: "device map contains unknown object"}
			return false
		}

		devices = append(devices, device)
		return true
	})

	return devices, err
}

func (r *SignatureDeviceInMemoryRepo) GetById(id string) (*domain.SignatureDevice, error) {
	item, ok := r.devices.Load(id)
	if !ok {
		return nil, errors.NotFoundError{Id: id}
	}

	device, ok := item.(*domain.SignatureDevice)
	if !ok {
		return nil, errors.InternalError{Msg: "device map contains unknown object"}
	}

	return device, nil
}

func (r *SignatureDeviceInMemoryRepo) UpdateSigningMetaInfo(id string, counter int, lastSignature string) (*domain.SignatureDevice, error) {
	device, err := r.GetById(id)
	if err != nil {
		return nil, err
	}

	device.SignatureCounter = counter
	device.LastSignature = lastSignature

	return device, nil
}
