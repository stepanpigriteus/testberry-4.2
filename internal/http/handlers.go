package httpsh

import (
	"encoding/json"
	"net/http"

	"disgreps/domain"
	"disgreps/internal/serv/worker"
)

// работяги мастер

func (s *Server) HandleReqOn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("reqon-ok"))
}

func (s *Server) HandleDone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("done-ok"))
}

// работяги

func (s *Server) HandleOn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("on-ok"))
}

func (s *Server) HandleLoad(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	type WorkerRequest struct {
		Chunk []domain.Line `json:"chunk"`
		Cfg   domain.Config `json:"cfg"`
	}

	var req WorkerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	result := worker.Worker(req.Cfg, req.Chunk)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
