package domain

import (
	"errors"
)

// Queue represents a virtual queue for a business.
type Queue struct {
	BusinessID string
	QueueID    string
	Users      []string
}

// QueueStatus represents the public state of the queue to be published.
type QueueStatus struct {
	BusinessID string   `json:"businessId"`
	QueueID    string   `json:"queueId"`
	Users      []string `json:"users"`
	Count      int      `json:"count"`
}

// NewQueue creates a new empty queue.
func NewQueue(businessID, queueID string) *Queue {
	return &Queue{
		BusinessID: businessID,
		QueueID:    queueID,
		Users:      make([]string, 0),
	}
}

// CanJoin checks if a user can join the queue.
func (q *Queue) CanJoin(userID string) error {
	if userID == "" {
		return errors.New("userID cannot be empty")
	}
	for _, u := range q.Users {
		if u == userID {
			return errors.New("user already in queue")
		}
	}
	return nil
}

// AddUser adds a user to the queue.
// It assumes CanJoin has been called and returned nil.
func (q *Queue) AddUser(userID string) {
	q.Users = append(q.Users, userID)
}

// ToStatus converts the Queue entity to its status representation.
func (q *Queue) ToStatus() QueueStatus {
	return QueueStatus{
		BusinessID: q.BusinessID,
		QueueID:    q.QueueID,
		Users:      q.Users,
		Count:      len(q.Users),
	}
}
