package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.temporal.io/sdk/client"

	"red-duck/internal/core/domain"
	"red-duck/internal/pkg/update"
)

type QueueHandler struct {
	Client client.Client
}

func (h *QueueHandler) JoinQueue(w http.ResponseWriter, r *http.Request) {
	businessID := r.URL.Query().Get("business_id")
	queueID := r.URL.Query().Get("queue_id")

	if businessID == "" || queueID == "" {
		http.Error(w, "missing business_id or queue_id", http.StatusBadRequest)
		return
	}

	var req domain.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	workflowID := fmt.Sprintf("%s:%s", businessID, queueID)
	joinQueueUpdate := update.New[domain.JoinRequest, int]("JoinQueue")

	// Call UpdateWorkflow and wait for it to complete
	handle, err := h.Client.UpdateWorkflow(r.Context(), client.UpdateWorkflowOptions{
		WorkflowID:   workflowID,
		UpdateID:     fmt.Sprintf("join-%s-%d", req.UserID, time.Now().UnixNano()),
		WaitForStage: client.WorkflowUpdateStageCompleted,
		UpdateName:   joinQueueUpdate.Name(),
		Args:         []interface{}{req},
	})

	if err != nil {
		// If the validator rejects the update, UpdateWorkflow returns an error.
		// We return 409 Conflict as requested.
		http.Error(w, fmt.Sprintf("Update rejected or failed: %v", err), http.StatusConflict)
		return
	}

	var position int
	if err := handle.Get(r.Context(), &position); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get update result: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]int{"position": position}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
