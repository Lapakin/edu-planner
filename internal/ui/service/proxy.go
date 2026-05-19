package service

import (
	"net/http"
)

func (s *Services) ProxyAuth(w http.ResponseWriter, r *http.Request) {
	s.AuthProxy.ServeHTTP(w, r)
}

func (s *Services) ProxySyllabus(w http.ResponseWriter, r *http.Request) {
	s.SyllabusProxy.ServeHTTP(w, r)
}
