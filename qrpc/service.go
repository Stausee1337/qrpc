package qrpc

import (
	"context"
	"fmt"
	"net/http"
)

type OperationKind uint

const (
	OpQuery OperationKind = iota
	OpMutation
)

type Service struct {
	name        string
	operations  [OpMutation + 1]*OperationResolver
	resolver 	resolver
	middlewares []func(http.Handler) http.Handler
}

func NewService(name string, resolvers ...*OperationResolver) *Service {
	operations := [OpMutation + 1]*OperationResolver{}
	for _, resolver := range resolvers {
		slot := &operations[resolver.kind]
		if *slot != nil {
			panic(fmt.Sprintf("qrpc.NewService: received multiple resolver assignments for %v", resolver.kind))
		}
		*slot = resolver
	}
	
	s := &Service{
		name: name,
		operations: operations,
	}
	s.resolver = s.buildHandleResolver()
	return s
}

func (s *Service) Use(middlewares ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middlewares...)
	s.resolver = s.buildHandleResolver()
}

func (s *Service) buildHandleResolver() resolver {
	commonMiddlewareHandler := chain(s.middlewares, http.HandlerFunc(callContextHandler))
	return handleResolver(func(kind OperationKind, op string) http.Handler {
		resolver := s.operations[kind]
		if resolver == nil {
			return nil
		}
		handler := resolver.handlersClosure(op)
		if handler == nil {
			return nil
		}
		return wrapWithContextHandler(commonMiddlewareHandler, handler)
	})
}

func chain(middlewares []func(http.Handler) http.Handler, endpoint http.Handler) http.Handler {
	if len(middlewares) == 0 {
		return endpoint
	}

	middleware := middlewares[len(middlewares) - 1](endpoint)
	for idx := len(middlewares) - 2; idx >= 0; idx-- {
		middleware = middlewares[idx](middleware)
	}

	return middleware
}


type contextKey struct{}

var wrappedHandlerKey = contextKey{}

func callContextHandler(w http.ResponseWriter, req *http.Request) {
	finalHandler := req.Context().Value(wrappedHandlerKey).(http.Handler)
	finalHandler.ServeHTTP(w, req)
}

func wrapWithContextHandler(baseHandler http.Handler, finalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newContext := context.WithValue(r.Context(), wrappedHandlerKey, finalHandler)
		baseHandler.ServeHTTP(w, r.WithContext(newContext))
	})
}

type handleResolver func(kind OperationKind, op string) http.Handler

func (r handleResolver) resolve(kind OperationKind, op string) http.Handler {
	return r(kind, op)
}

type OperationResolver struct {
	kind OperationKind
	handlersClosure func(op string) http.Handler
}

func BindResolver[H any](
	kind OperationKind,
	handlers H,
	resolver func(handlers H, op string) http.Handler,
) *OperationResolver {
	return &OperationResolver{
		kind: kind,
		handlersClosure: func(op string) http.Handler {
			return resolver(handlers, op)
		},
	}
}

func (r *OperationResolver) resolve(kind OperationKind, op string) http.Handler {
	if kind != r.kind {
		return nil;
	}
	return r.handlersClosure(op)
}



