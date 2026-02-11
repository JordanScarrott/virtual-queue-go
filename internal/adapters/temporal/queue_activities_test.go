package temporal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventTracker struct
type MockEventTracker struct {
	mock.Mock
}

func (m *MockEventTracker) Track(eventType, businessID, userID string, props map[string]interface{}) error {
	args := m.Called(eventType, businessID, userID, props)
	return args.Error(0)
}

func TestJoinQueue_TracksEvent(t *testing.T) {
	mockTracker := new(MockEventTracker)

	activities := &QueueActivities{
		Tracker: mockTracker,
	}

	params := JoinQueueParams{
		BusinessID:      "biz_123",
		UserID:          "user_456",
		QueueLength:     5,
		WaitTimeMinutes: 10,
	}

	// Expect Track to be called
	mockTracker.On("Track", "queue.joined", "biz_123", "user_456", mock.Anything).Return(nil)

	err := activities.JoinQueue(context.Background(), params)

	assert.NoError(t, err)
	mockTracker.AssertExpectations(t)
}
