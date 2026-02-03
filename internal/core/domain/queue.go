package domain

// JoinRequest represents the payload to join a queue.
type JoinRequest struct {
	UserID string `json:"user_id"`
}

// Queue represents the state of a business queue.
type Queue struct {
	BusinessID string
	QueueID    string
	Closed     bool
	Users      []string // List of UserIDs in order
}

// NewQueue creates a new Queue.
func NewQueue(businessID, queueID string) *Queue {
	return &Queue{
		BusinessID: businessID,
		QueueID:    queueID,
		Closed:     false,
		Users:      make([]string, 0),
	}
}
