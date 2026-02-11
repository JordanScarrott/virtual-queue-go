package domain

import (
	"errors"
	"time"
)

var (
	ErrUserAlreadyInQueue = errors.New("user already in queue")
	ErrUserNotFound       = errors.New("user not found in queue")
)

type TicketStatus string

const (
	TicketStatusWaiting   TicketStatus = "WAITING"
	TicketStatusReady     TicketStatus = "READY"
	TicketStatusCompleted TicketStatus = "COMPLETED"
)

type Ticket struct {
	UserID     string       `json:"userId"`
	Status     TicketStatus `json:"status"`
	AssignedTo string       `json:"assignedTo,omitempty"` // Counter ID, e.g. "Counter 3"
	JoinedAt   time.Time    `json:"joinedAt"`
}

type Queue struct {
	ID         string
	BusinessID string
	Tickets    []Ticket
}

type JoinRequest struct {
	UserID string `json:"userId"`
}

func NewQueue(id, businessID string) *Queue {
	return &Queue{
		ID:         id,
		BusinessID: businessID,
		Tickets:    make([]Ticket, 0),
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
	for _, t := range q.Tickets {
		if t.UserID == userID {
			return ErrUserAlreadyInQueue
		}
	}
	return nil
}

func (q *Queue) AddUser(userID string) int {
	ticket := Ticket{
		UserID:   userID,
		Status:   TicketStatusWaiting,
		JoinedAt: time.Now(),
	}
	q.Tickets = append(q.Tickets, ticket)
	return len(q.Tickets)
}

// Dequeue removes a user from the queue by ID.
func (q *Queue) Dequeue(userID string) error {
	for i, t := range q.Tickets {
		if t.UserID == userID {
			q.Tickets = append(q.Tickets[:i], q.Tickets[i+1:]...)
			return nil
		}
	}
	return ErrUserNotFound
}

// GetPosition returns the 1-based index of the user in the queue.
// Returns 0 if not found.
func (q *Queue) GetPosition(userID string) int {
	for i, t := range q.Tickets {
		if t.UserID == userID {
			return i + 1
		}
	}
	return 0
}

// Len returns the number of users in the queue.
func (q *Queue) Len() int {
	return len(q.Tickets)
}

// ServeNext finds the next waiting ticket, updates its status to READY and assigns it to the counter.
func (q *Queue) ServeNext(counterID string) (*Ticket, error) {
	for i := range q.Tickets {
		if q.Tickets[i].Status == TicketStatusWaiting {
			q.Tickets[i].Status = TicketStatusReady
			q.Tickets[i].AssignedTo = counterID
			return &q.Tickets[i], nil
		}
	}
	return nil, errors.New("queue empty")
}

// Snapshot returns a copy of the current state
func (q *Queue) Snapshot() Queue {
	ticketsCopy := make([]Ticket, len(q.Tickets))
	copy(ticketsCopy, q.Tickets)
	return Queue{
		ID:         q.ID,
		BusinessID: q.BusinessID,
		Tickets:    ticketsCopy,
	}
}
