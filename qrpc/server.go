package qrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dlclark/regexp2"
)

type Server struct {
	logger *slog.Logger
	services []*Service
	resolvers map[string]resolver
}

func NewRPCServer(services ...*Service) *Server {
	resolvers := map[string]resolver{}
	for _, service := range services {
		resolvers[service.name] = service.resolver
	}

	return &Server{
		services: services,
		resolvers: resolvers,
	}
}

func (s *Server) SetLogger(logger *slog.Logger) {
	s.logger = logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed)
		return;
	}

	defer func() {
		if err := recover(); err != nil {
			if s.logger != nil {
				aerr, ok := err.(error)
				if ok {
					s.logger.Error(aerr.Error())
				} else {
					s.logger.Error(fmt.Sprintf("%v", err))
				}
			}
			writeError(w, http.StatusInternalServerError)
		}
	}()


	params := parseParams(r)
	if params == nil {
		writeError(w, http.StatusBadRequest)
		return
	}

	handler := s.resolveRequest(params)
	if handler == nil {
		writeError(w, http.StatusNotFound)
		return
	}

	handler.ServeHTTP(w, r)
}

func (s *Server) resolveRequest(request *operationParams) http.Handler {
	serviceResolver, ok := s.resolvers[request.service]
	if !ok {
		return nil;
	}
	
	return serviceResolver.resolve(request.kind, request.op);
}

type operationParams struct {
	service string
	kind    OperationKind
	op      string
}

var opRegex = regexp2.MustCompile(`(?<kind>query|mutation):(?<service>[A-Za-z][A-Za-z0-9]*)\.(?<operation>[A-Za-z][A-Za-z0-9]*)`, regexp2.Compiled)

func parseParams(r *http.Request) *operationParams {
	urlPath := r.URL.Path

	segments := strings.Split(strings.Trim(urlPath, "/"), "/")
	querySegment := segments[len(segments) - 1]

	match, err := opRegex.FindStringMatch(querySegment)
	if err != nil || match == nil {
		return nil;
	}

	kindGroup := match.GroupByName("kind")
	serviceGroup := match.GroupByName("service")
	opGroup := match.GroupByName("operation")

	var kind OperationKind
	switch kindGroup.String() {
	case "query":
		kind = OpQuery
	case "mutation":
		kind = OpMutation
	default:
		return nil;
	}

	service := serviceGroup.String()
	op := opGroup.String()

	return &operationParams{ service, kind, op }
}

func unmarshalRequestHelper(r *http.Request, a any) error {
	body, err := io.ReadAll(r.Body);
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, a)

	switch err.(type) {
	case *json.SyntaxError, *json.UnmarshalTypeError:
		return err
	case nil:
		return nil
	default:
		panic(err)
	}
}

func MakeHandler(handler func(ctx context.Context) Response) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := handler(r.Context())

		data := response.rawData()
		message, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		w.Write(message)
	})
}

func MakeInputHandler[I any](handler func(ctx context.Context, input I) Response) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var input I

		err := unmarshalRequestHelper(r, &input)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest) + ": " + err.Error(), http.StatusBadRequest)
			return
		}

		response := handler(r.Context(), input)
		writeResponse(w, response)
	})
}

func writeError(w http.ResponseWriter, statusCode int) {
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func writeResponse(w http.ResponseWriter, response Response) {
	var data []byte
	if response != nil {
		var err error
		data, err = json.Marshal(response.rawData())
		if err != nil {
			panic(err)
		}
	} else {
		data = []byte("{}")
	}

	w.WriteHeader(200)
	w.Write(data)
}
