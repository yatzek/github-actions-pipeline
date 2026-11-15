package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"gist-service/client"
)

type Server struct {
	client client.GistClient
	mux    *http.ServeMux
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewServer(client client.GistClient) *Server {
	s := &Server{
		client: client,
		mux:    http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /{$}", s.handleRoot)
	s.mux.HandleFunc("GET /{user}", s.handleGetUserGists)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	defer func() {
		if rec := recover(); rec != nil {
			log.Printf("panic recovered: %v", rec)
			s.writeJSONError(w, http.StatusInternalServerError, "internal server error")
		}
		log.Printf("[HTTP] %s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	}()

	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"GitHub Gists API. Query /<username> to retrieve public gists."}`))
}

func (s *Server) handleGetUserGists(w http.ResponseWriter, r *http.Request) {

	user := strings.TrimSpace(r.PathValue("user"))
	if user == "" {
		s.writeJSONError(w, http.StatusBadRequest, "username parameter is required")
		return
	}

	gists, resp, err := s.client.ListUserGists(r.Context(), user)

	if err != nil {
		log.Printf("error fetching gists for user %q: %v", user, err)
		s.writeJSONError(w, resp.StatusCode, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(gists); err != nil {
		log.Printf("encoding gists response: %v", err)
	}
}

func (s *Server) writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	resp := errorResponse{
		Error: message,
	}
	_ = json.NewEncoder(w).Encode(resp)
}
