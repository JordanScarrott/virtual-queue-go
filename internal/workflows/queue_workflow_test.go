package workflows

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

	"example.com/virtual-queue/internal/core/domain"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) TestQueueWorkflow() {
	// Schedule signals
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalJoinQueue, JoinQueueSignal{UserID: "u1"})
	}, time.Second)

	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalJoinQueue, JoinQueueSignal{UserID: "u2"})
	}, time.Second*2)

	// Check state after joins
	s.env.RegisterDelayedCallback(func() {
		res, err := s.env.QueryWorkflow(QueryGetState)
		s.NoError(err)
		var state domain.Queue
		err = res.Get(&state)
		s.NoError(err)
		s.Equal(2, len(state.Users))
		s.Equal("u1", state.Users[0])
		s.Equal("u2", state.Users[1])
	}, time.Second*3)

	// Leave
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalLeaveQueue, LeaveQueueSignal{UserID: "u1"})
	}, time.Second*4)

	// Check state after leave
	s.env.RegisterDelayedCallback(func() {
		res, err := s.env.QueryWorkflow(QueryGetState)
		s.NoError(err)
		var state domain.Queue
		err = res.Get(&state)
		s.NoError(err)
		s.Equal(1, len(state.Users))
		s.Equal("u2", state.Users[0])
	}, time.Second*5)

	// Exit
	s.env.RegisterDelayedCallback(func() {
		s.env.SignalWorkflow(SignalExit, nil)
	}, time.Second*6)

	// Execute Workflow
	s.env.ExecuteWorkflow(QueueWorkflow)

	s.True(s.env.IsWorkflowCompleted())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
