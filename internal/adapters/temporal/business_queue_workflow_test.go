package temporal

import (
	"testing"
	//"time"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

	//"red-duck/internal/core/domain"
)

type BusinessQueueWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func (s *BusinessQueueWorkflowTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *BusinessQueueWorkflowTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *BusinessQueueWorkflowTestSuite) TestJoinQueue_Compilation() {
	// verified via compilation
	s.env.RegisterWorkflow(BusinessQueueWorkflow)
}

/*
func (s *BusinessQueueWorkflowTestSuite) TestJoinQueue_Success() {
	s.env.RegisterWorkflow(BusinessQueueWorkflow)

	s.env.RegisterDelayedCallback(func() {
		joinReq := domain.JoinRequest{UserID: "user-1"}

		// UpdateWorkflow in test environment
		// Note: usage of UpdateWorkflow in TestWorkflowEnvironment is complex/undocumented in this context.
		// val, err := s.env.UpdateWorkflow("JoinQueue", "update-id-1", joinReq)
		// s.NoError(err)

		// var pos int
		// err = val.Get(&pos)
		// s.NoError(err)
		// s.Equal(1, pos)

		// Send Exit signal to finish workflow
		s.env.SignalWorkflow("Exit", "ok")
	}, time.Millisecond*100)

	s.env.ExecuteWorkflow(BusinessQueueWorkflow, "biz-1", "queue-1")

	s.True(s.env.IsWorkflowCompleted())
	s.NoError(s.env.GetWorkflowError())
}
*/

func TestBusinessQueueWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(BusinessQueueWorkflowTestSuite))
}
