package domain

import "errors"

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

// CanJoin checks if a user can join the queue.
// It returns an error if the queue is closed or the user is already in the queue.
func (q *Queue) CanJoin(userID string) error {
	if q.Closed {
		return errors.New("queue is closed")
	}
	for _, uid := range q.Users {
		if uid == userID {
			return errors.New("user already in queue")
		}
	}
	return nil
}

// AddUser adds a user to the queue and returns their position (1-based).
// It assumes CanJoin has already been called or validation is handled by the caller.
// Ideally, CanJoin should be called before this.
func (q *Queue) AddUser(userID string) int {
	q.Users = append(q.Users, userID)
	return len(q.Users)
}
