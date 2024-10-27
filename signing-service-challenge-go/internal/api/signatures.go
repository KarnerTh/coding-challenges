package api

import (
	"encoding/json"
	"net/http"
)

type createSignatureRequest struct {
	Data string `json:"data" validate:"required"`
}

func (s *Server) Signatures(response http.ResponseWriter, request *http.Request) {
	deviceId := request.PathValue("deviceId")
	if len(deviceId) == 0 {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"device id is missing in path",
		})
	}

	switch request.Method {
	case http.MethodPost:
		if !isMediaType(request, ApplicationJson) {
			WriteErrorResponse(response, http.StatusUnsupportedMediaType, nil)
			return
		}

		// Limit the request body to 1MB
		request.Body = http.MaxBytesReader(response, request.Body, 1<<20)

		var signatureRequest createSignatureRequest
		err := json.NewDecoder(request.Body).Decode(&signatureRequest)
		if err != nil {
			WriteErrorResponse(response, http.StatusBadRequest, []string{
				err.Error(),
			})
			return
		}

		err = validate.Struct(signatureRequest)
		if err != nil {
			WriteErrorResponse(response, http.StatusBadRequest, []string{
				err.Error(),
			})
			return
		}

		signature, err := s.signatureDeviceService.Sign(deviceId, signatureRequest.Data)
		if err != nil {
			WriteErrorResponse(response, statusCodeFromError(err), []string{
				err.Error(),
			})
			return
		}

		WriteAPIResponse(response, http.StatusCreated, signature)

	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil)
	}
}
