package domain

import (
	"testing"
)

func TestQueue_Enqueue(t *testing.T) {
	q := NewQueue("q1", "biz1")

	if q.BusinessID != "biz1" {
		t.Errorf("expected businessID biz1, got %s", q.BusinessID)
	}

	if err := q.Enqueue("u1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if q.Len() != 1 {
		t.Errorf("expected len 1, got %d", q.Len())
	}

	if err := q.Enqueue("u1"); err != ErrUserAlreadyInQueue {
		t.Errorf("expected ErrUserAlreadyInQueue, got %v", err)
	}
}

func TestQueue_Dequeue(t *testing.T) {
	q := NewQueue("q1", "biz1")
	q.Enqueue("u1")
	q.Enqueue("u2")

	if err := q.Dequeue("u1"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if q.Len() != 1 {
		t.Errorf("expected len 1, got %d", q.Len())
	}

	if q.Users[0] != "u2" {
		t.Errorf("expected u2 at head, got %s", q.Users[0])
	}

	if err := q.Dequeue("u3"); err != ErrUserNotFound {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}

func TestQueue_GetPosition(t *testing.T) {
	q := NewQueue("q1", "biz1")
	q.Enqueue("u1")
	q.Enqueue("u2")

	if pos := q.GetPosition("u1"); pos != 1 {
		t.Errorf("expected pos 1, got %d", pos)
	}
	if pos := q.GetPosition("u2"); pos != 2 {
		t.Errorf("expected pos 2, got %d", pos)
	}
	if pos := q.GetPosition("u3"); pos != 0 {
		t.Errorf("expected pos 0, got %d", pos)
	}
}
