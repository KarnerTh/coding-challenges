package api

import (
	"encoding/json"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress          string
	signatureDeviceService domain.SignatureDeviceService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string) *Server {
	return &Server{
		listenAddress:          listenAddress,
		signatureDeviceService: domain.NewSignatureDeviceService(persistence.NewSignatureDeviceInMemoryRepo()),
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))
	mux.Handle("/api/v0/devices", http.HandlerFunc(s.Devices))
	mux.Handle("/api/v0/devices/{deviceId}", http.HandlerFunc(s.Device))
	mux.Handle("/api/v0/devices/{deviceId}/signatures", http.HandlerFunc(s.Signatures))

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: append([]string{http.StatusText(code)}, errors...),
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

func statusCodeFromError(error error) int {
	switch error.(type) {
	case errors.NotFoundError:
		return http.StatusNotFound
	case errors.ConflictError:
		return http.StatusConflict
	case errors.BadInputError:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
