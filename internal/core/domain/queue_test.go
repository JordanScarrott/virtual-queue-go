package domain

import (
	"testing"
)

func TestQueue_CanJoin(t *testing.T) {
	q := NewQueue("biz1", "q1")

	// Test adding a user
	err := q.CanJoin("user1")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	q.AddUser("user1")

	// Test adding same user again
	err = q.CanJoin("user1")
	if err == nil {
		t.Error("Expected error for duplicate user, got nil")
	}

	// Test adding another user
	err = q.CanJoin("user2")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
}
