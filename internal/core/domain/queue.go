package domain

import (
	"errors"
	"time"
)

var (
	ErrUserAlreadyInQueue = errors.New("user already in queue")
	ErrUserNotInQueue     = errors.New("user not in queue")
)

type QueueItem struct {
	UserID    string
	JoinedAt  time.Time
}

type Queue struct {
	ID    string
	Items []QueueItem
}

func NewQueue(id string) *Queue {
	return &Queue{
		ID:    id,
		Items: make([]QueueItem, 0),
	}
}

func (q *Queue) Join(userID string, now time.Time) error {
	for _, item := range q.Items {
		if item.UserID == userID {
			return ErrUserAlreadyInQueue
		}
	}
	q.Items = append(q.Items, QueueItem{
		UserID:   userID,
		JoinedAt: now,
	})
	return nil
}

func (q *Queue) Leave(userID string) error {
	for i, item := range q.Items {
		if item.UserID == userID {
			q.Items = append(q.Items[:i], q.Items[i+1:]...)
			return nil
		}
	}
	return ErrUserNotInQueue
}

func (q *Queue) Position(userID string) (int, error) {
	for i, item := range q.Items {
		if item.UserID == userID {
			return i + 1, nil // 1-based index
		}
	}
	return 0, ErrUserNotInQueue
}
