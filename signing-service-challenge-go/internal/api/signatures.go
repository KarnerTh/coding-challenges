package api

import "net/http"

func (s *Server) Signatures(response http.ResponseWriter, request *http.Request) {
	deviceId := request.PathValue("deviceId")
	if len(deviceId) == 0 {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			"device id is missing in path",
		})
	}

	switch request.Method {
	case http.MethodPost:
	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil)
	}
}
