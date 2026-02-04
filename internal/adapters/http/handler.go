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
	Client    client.Client
	TaskQueue string
}

func (h *QueueHandler) CreateQueue(w http.ResponseWriter, r *http.Request) {
	businessID := r.URL.Query().Get("business_id")
	queueID := r.URL.Query().Get("queue_id")

	if businessID == "" || queueID == "" {
		http.Error(w, "missing business_id or queue_id", http.StatusBadRequest)
		return
	}

	workflowID := fmt.Sprintf("%s:%s", businessID, queueID)

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: h.TaskQueue,
	}

	// We use the string name "BusinessQueueWorkflow" to avoid importing the adapter package
	run, err := h.Client.ExecuteWorkflow(r.Context(), options, "BusinessQueueWorkflow", businessID, queueID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to start workflow: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"workflow_id": run.GetID(),
		"run_id":      run.GetRunID(),
	})
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

func (h *QueueHandler) LeaveQueue(w http.ResponseWriter, r *http.Request) {
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
	leaveQueueUpdate := update.New[domain.JoinRequest, int]("LeaveQueue")

	handle, err := h.Client.UpdateWorkflow(r.Context(), client.UpdateWorkflowOptions{
		WorkflowID:   workflowID,
		UpdateID:     fmt.Sprintf("leave-%s-%d", req.UserID, time.Now().UnixNano()),
		WaitForStage: client.WorkflowUpdateStageCompleted,
		UpdateName:   leaveQueueUpdate.Name(),
		Args:         []interface{}{req},
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Update failed: %v", err), http.StatusInternalServerError)
		return
	}

	var remaining int
	if err := handle.Get(r.Context(), &remaining); err != nil {
		http.Error(w, fmt.Sprintf("Failed to get update result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"remaining_users": remaining})
}

func (h *QueueHandler) GetQueueStatus(w http.ResponseWriter, r *http.Request) {
	businessID := r.URL.Query().Get("business_id")
	queueID := r.URL.Query().Get("queue_id")

	if businessID == "" || queueID == "" {
		http.Error(w, "missing business_id or queue_id", http.StatusBadRequest)
		return
	}

	workflowID := fmt.Sprintf("%s:%s", businessID, queueID)

	response, err := h.Client.QueryWorkflow(r.Context(), workflowID, "", "GetStatus")
	if err != nil {
		http.Error(w, fmt.Sprintf("Query failed: %v", err), http.StatusInternalServerError)
		return
	}

	var q domain.Queue
	if err := response.Get(&q); err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode query result: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(q)
}
