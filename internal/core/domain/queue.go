package domain

import (
	"errors"
)

var (
	ErrUserAlreadyInQueue = errors.New("user already in queue")
	ErrUserNotFound       = errors.New("user not found in queue")
)

type Queue struct {
	ID         string
	BusinessID string
	Users      []string
}

type JoinRequest struct {
	UserID string `json:"userId"`
}

func NewQueue(id, businessID string) *Queue {
	return &Queue{
		ID:         id,
		BusinessID: businessID,
		Users:      make([]string, 0),
	}
}

// Enqueue adds a user to the end of the queue.
func (q *Queue) Enqueue(userID string) error {
	if err := q.CanJoin(userID); err != nil {
		return err
	}
	q.AddUser(userID)
	return nil
}

func (q *Queue) CanJoin(userID string) error {
	for _, u := range q.Users {
		if u == userID {
			return ErrUserAlreadyInQueue
		}
	}
	return nil
}

func (q *Queue) AddUser(userID string) int {
	q.Users = append(q.Users, userID)
	return len(q.Users)
}

// Dequeue removes a user from the queue by ID.
func (q *Queue) Dequeue(userID string) error {
	for i, u := range q.Users {
		if u == userID {
			q.Users = append(q.Users[:i], q.Users[i+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}

// GetPosition returns the 1-based index of the user in the queue.
// Returns 0 if not found.
func (q *Queue) GetPosition(userID string) int {
	for i, u := range q.Users {
		if u == userID {
			return i + 1
		}
	}
	return 0
}

// Len returns the number of users in the queue.
func (q *Queue) Len() int {
	return len(q.Users)
}

// Snapshot returns a copy of the current state
func (q *Queue) Snapshot() Queue {
	usersCopy := make([]string, len(q.Users))
	copy(usersCopy, q.Users)
	return Queue{
		ID:         q.ID,
		BusinessID: q.BusinessID,
		Users:      usersCopy,
	}
}
