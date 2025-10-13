package qrpc

import (
	"net/http"
)

type Response interface {
	rawData() map[string]any
}

type successResponse struct { data any }
func (r *successResponse)	rawData() map[string]any {
	return map[string]any{
		"data": r.data,
	};
}

func WrapData(data any) Response {
	return &successResponse{ data: data }
}

type errorResponse struct { error error }
func (r *errorResponse)	rawData() map[string]any {
	return map[string]any{
		"error": r.error.Error(),
	};
}

func ErrorData(error error) Response {
	return &errorResponse{ error: error }
}

type resolver interface {
	resolve(kind OperationKind, op string) http.Handler
}

