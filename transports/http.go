package transports

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"

	"github.com/victordinh/gokit-grpc/endpoints"
	"github.com/victordinh/gokit-grpc/util"
)

func NewHTTPServer(ep endpoints.Endpoints, logger log.Logger) *mux.Router {
	//options provided by the Go kit to facilitate error control
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	//definition of a handler
	getHealthCheckHandler := httptransport.NewServer(
		ep.HealthCheckEndpoint, //use the endpoint
		func(context.Context, *http.Request) (interface{}, error) { return struct{}{}, nil }, //converts the parameters received via the request body into the struct expected by the endpoint
		encodeResponse, //converts the struct returned by the endpoint to a json response
		options...,
	)

	addHandler := httptransport.NewServer(
		ep.Add,           //use the endpoint
		decodeAddRequest, //converts the parameters received via the request body into the struct expected by the endpoint
		encodeResponse,   //converts the struct returned by the endpoint to a json response
		options...,
	)

	validateUserHandler := httptransport.NewServer(
		ep.ValidateUserEndpoint,   //use the endpoint
		decodeValidateUserRequest, //converts the parameters received via the request body into the struct expected by the endpoint
		encodeResponse,            //converts the struct returned by the endpoint to a json response
		options...,
	)

	validateTokenHandler := httptransport.NewServer(
		ep.ValidateTokenEndpoint,
		decodeValidateTokenRequest,
		encodeResponse,
		options...,
	)

	r := mux.NewRouter()
	r.Methods("GET").Path("/hc").Handler(getHealthCheckHandler)
	r.Methods("POST").Path("/v1/auth").Handler(validateUserHandler)
	r.Methods("POST").Path("/v1/validate-token").Handler(validateTokenHandler)
	r.Methods("POST").Path("/v1/add").Handler(addHandler)

	return r
}

func decodeAddRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.MathReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeValidateUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ValidateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeValidateTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.ValidateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := response.(error); ok && e != nil {
		encodeError(ctx, e, w)
		return nil
	}
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case util.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case util.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
