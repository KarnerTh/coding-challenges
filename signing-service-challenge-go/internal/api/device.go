package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
)

type createDeviceRequest struct {
	Id        string                    `json:"id" validate:"required"`
	Algorithm crypto.SignatureAlgorithm `json:"algorithm" validate:"required"`
	Label     string                    `json:"label"`
}

func (s *Server) Devices(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		devices, err := s.signatureDeviceService.GetAll()
		if err != nil {
			WriteErrorResponse(response, statusCodeFromError(err), []string{
				err.Error(),
			})
			return
		}

		WriteAPIResponse(response, http.StatusOK, devices)
	case http.MethodPost:
		if !isMediaType(request, ApplicationJson) {
			WriteErrorResponse(response, http.StatusUnsupportedMediaType, nil)
			return
		}

		// Limit the request body to 1MB
		request.Body = http.MaxBytesReader(response, request.Body, 1<<20)

		var deviceRequest createDeviceRequest
		err := json.NewDecoder(request.Body).Decode(&deviceRequest)
		if err != nil {
			WriteErrorResponse(response, http.StatusBadRequest, []string{
				err.Error(),
			})
			return
		}

		err = validate.Struct(deviceRequest)
		if err != nil {
			WriteErrorResponse(response, http.StatusBadRequest, []string{
				err.Error(),
			})
			return
		}

		device, err := s.signatureDeviceService.Create(
			deviceRequest.Id,
			deviceRequest.Algorithm,
			deviceRequest.Label,
		)
		if err != nil {
			WriteErrorResponse(response, statusCodeFromError(err), []string{
				err.Error(),
			})
			return
		}

		WriteAPIResponse(response, http.StatusCreated, device)

	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil)
	}
}

func (s *Server) Device(response http.ResponseWriter, request *http.Request) {
	deviceId := request.PathValue("deviceId")
	if len(deviceId) == 0 {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"deviceId is missing in path",
		})
		return
	}

	switch request.Method {
	case http.MethodGet:
		device, err := s.signatureDeviceService.GetById(deviceId)
		if err != nil {
			WriteErrorResponse(response, statusCodeFromError(err), []string{
				err.Error(),
			})
			return
		}

		WriteAPIResponse(response, http.StatusOK, device)

	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil)
	}
}
