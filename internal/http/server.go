package httpsh

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"disgreps/domain"
)

type Server struct {
	port string
	host string
	mode bool
	cfg  domain.Config
	srv  *http.Server
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewServer(port, host string, mode bool, cfg domain.Config) *Server {
	return &Server{
		port: port,
		host: host,
		mode: mode,
		cfg:  cfg,
	}
}

func (s *Server) RunServer(ctx context.Context) error {
	if s.port == "" || s.host == "" {
		log.Fatal("port or host empty")
	}

	mux := http.NewServeMux()
	mux.Handle("/", &handleDef{})

	if s.mode == true {
		mux.HandleFunc("/reqon", s.HandleReqOn)
		mux.HandleFunc("/done", s.HandleDone)
	} else {
		mux.HandleFunc("/on", s.HandleOn)
		mux.HandleFunc("/load", s.HandleLoad)
	}

	srv := &http.Server{
		Addr:         s.host + ":" + s.port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	s.srv = srv
	fmt.Println("сервер стартовал)")
	return s.srv.ListenAndServe()
}

type handleDef struct{}

func (h *handleDef) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	statusCode := http.StatusNotFound
	if r.Method == http.MethodOptions {
		statusCode = http.StatusOK
	}

	w.WriteHeader(statusCode)
	response := ErrorResponse{
		Message: "Endpoint not found or method not allowed",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
	}
}
