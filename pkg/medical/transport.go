package medical

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	ErrBadRouting = errors.New("route and handle error")
)

func MakeHTTpHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	endpoint := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/medical/").Handler(httptransport.NewServer(
		endpoint.PostVisitEndpoints,
		decodePostVisitRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/medical/{id}").Handler(httptransport.NewServer(
		endpoint.GetVisitEndpoints,
		decodeGetVisitRequest,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/patient/").Handler(httptransport.NewServer(
		endpoint.PostPatientEndpoints,
		decodePostPatientRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/patient/{id}").Handler(httptransport.NewServer(
		endpoint.GetPatientEndpoints,
		decodeGetPatientRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostVisitRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postVisit
	if e := json.NewDecoder(r.Body).Decode(&req.Visit); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetVisitRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getVisit{ID: id}, nil
}

func decodePostPatientRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postPatient
	if e := json.NewDecoder(r.Body).Decode(&req.Patient); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetPatientRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getPatient{ID: id}, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
