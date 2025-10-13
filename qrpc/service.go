package qrpc

import (
	"encoding/json"
	"net/http"
)

type operation struct {
}

type ServiceBuilder struct {
	name string
}

func CreateBuilderForService(name string) *ServiceBuilder {
	return &ServiceBuilder{
		name,
	}
}

type Service struct {
	name        string
	mutations   map[string]operation
	queries     map[string]operation
	middlewares []func(http.Handler) http.Handler
}

func (s *Service) Use(middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
}
