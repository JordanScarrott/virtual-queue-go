package analytics

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) InsertEvent(ctx context.Context, payload EventPayload) error {
	args := m.Called(ctx, payload)
	return args.Error(0)
}

func TestProcessMessage_InsertsToDB(t *testing.T) {
	mockRepo := new(MockEventRepository)

	payload := EventPayload{
		Type:       "queue.joined",
		BusinessID: "biz_123",
		UserID:     "123",
		Timestamp:  time.Now(),
		Properties: map[string]interface{}{"foo": "bar"},
	}

	jsonBytes, err := json.Marshal(payload)
	assert.NoError(t, err)

	// Expect InsertEvent to be called
	// We verify that the payload passed to InsertEvent has the correct UserID
	mockRepo.On("InsertEvent", mock.Anything, mock.MatchedBy(func(p EventPayload) bool {
		return p.UserID == "123" && p.Type == "queue.joined"
	})).Return(nil)

	err = ProcessMessage(jsonBytes, mockRepo)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
