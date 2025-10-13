package qrpc

import "net/http"

type Server struct {}

func NewRPCServer(services ...Service) *Server {
	return &Server{}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {

}

