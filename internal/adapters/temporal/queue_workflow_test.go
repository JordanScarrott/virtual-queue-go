package temporal

import (
	"testing"
	"time"

	"red-duck/internal/core/domain"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type QueueWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *QueueWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *QueueWorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *QueueWorkflowTestSuite) TestJoinAndQuery() {
	queueID := "test-queue"
	userID := "user-1"

	s.env.RegisterWorkflow(QueueWorkflow)

	// Step 1: Join Queue
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalJoinQueue, JoinQueueSignal{UserID: userID})
	}, time.Millisecond*10)

	// Step 2: Query Queue to verify Join
	s.env.RegisterDelayedCallback(func() {
		val, err := s.env.QueryWorkflow(QueryGetQueue)
		s.NoError(err)
		var q domain.Queue
		err = val.Get(&q)
		s.NoError(err)
		s.Equal(1, len(q.Items))
		s.Equal(userID, q.Items[0].UserID)

		// Step 3: Leave Queue
		s.env.SignalWorkflow(SignalLeaveQueue, LeaveQueueSignal{UserID: userID})
	}, time.Millisecond*20)

	// Step 4: Query Queue to verify Leave
	s.env.RegisterDelayedCallback(func() {
		val, err := s.env.QueryWorkflow(QueryGetQueue)
		s.NoError(err)
		var q domain.Queue
		err = val.Get(&q)
		s.NoError(err)
		s.Equal(0, len(q.Items))

		// Cancel to finish workflow
		s.env.CancelWorkflow()
	}, time.Millisecond*30)

	s.env.ExecuteWorkflow(QueueWorkflow, queueID)

	s.True(s.env.IsWorkflowCompleted())
	s.Error(s.env.GetWorkflowError())
}

func TestQueueWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(QueueWorkflowTestSuite))
}
