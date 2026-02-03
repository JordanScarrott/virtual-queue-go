package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"red-duck/internal/core/ports"
)

type QueueHandler struct {
	Service ports.QueueService
}

func NewQueueHandler(service ports.QueueService) *QueueHandler {
	return &QueueHandler{Service: service}
}

type UserRequest struct {
	UserID string `json:"userID"`
}

func (h *QueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !strings.HasPrefix(path, "/queues/") {
		http.NotFound(w, r)
		return
	}

	parts := strings.Split(strings.TrimPrefix(path, "/queues/"), "/")
	queueID := parts[0]
	if queueID == "" {
		http.Error(w, "queue ID required", http.StatusBadRequest)
		return
	}

	if len(parts) == 1 {
		if r.Method == http.MethodPost {
			h.createQueue(w, r, queueID)
		} else if r.Method == http.MethodGet {
			h.getQueue(w, r, queueID)
		} else {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	if len(parts) == 2 {
		action := parts[1]
		if action == "join" && r.Method == http.MethodPost {
			h.joinQueue(w, r, queueID)
		} else if action == "leave" && r.Method == http.MethodPost {
			h.leaveQueue(w, r, queueID)
		} else {
			http.NotFound(w, r)
		}
		return
	}

	http.NotFound(w, r)
}

func (h *QueueHandler) createQueue(w http.ResponseWriter, r *http.Request, queueID string) {
	err := h.Service.CreateQueue(r.Context(), queueID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *QueueHandler) getQueue(w http.ResponseWriter, r *http.Request, queueID string) {
	q, err := h.Service.GetQueue(r.Context(), queueID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}

func (h *QueueHandler) joinQueue(w http.ResponseWriter, r *http.Request, queueID string) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, "userID required", http.StatusBadRequest)
		return
	}

	err := h.Service.JoinQueue(r.Context(), queueID, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *QueueHandler) leaveQueue(w http.ResponseWriter, r *http.Request, queueID string) {
	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, "userID required", http.StatusBadRequest)
		return
	}

	err := h.Service.LeaveQueue(r.Context(), queueID, req.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
