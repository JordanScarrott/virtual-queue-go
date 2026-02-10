package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) TestLoginWorkflow_Success() {
	env := s.NewTestWorkflowEnvironment()

	// Register Workflow
	env.RegisterWorkflow(LoginWorkflow)
	// Register Activities (Best practice even when mocking to ensure signatures match)
	env.RegisterActivity(SendMagicCode)
	env.RegisterActivity(GenerateToken)

	var generatedCode string

	// Mock SendMagicCode
	// Note: OnActivity matches arguments passed to ExecuteActivity.
	// SendMagicCode signature: func(ctx, email, code)
	// ExecuteActivity passes: email, code
	// So we match (mock.Anything, mock.Anything) for (email, code).
	// Context is typically ignored in OnActivity matching for struct payloads but activity functions might differ.
	// Let's use flexible matching and logging.
	env.OnActivity(SendMagicCode, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			// Debugging args
			for i, arg := range args {
				s.T().Logf("SendMagicCode Arg %d: %v", i, arg)
			}
			// args[0] is usually Context, args[1] is Email, args[2] is Code
			if len(args) > 2 {
				generatedCode = args.String(2)
			} else {
				s.T().Log("SendMagicCode mock received fewer args than expected")
			}
		})

	// Mock GenerateToken
	// GenerateToken signature: func(ctx, User)
	// ExecuteActivity passes: User
	// We match (mock.Anything, mock.Anything) for (ctx, User)
	env.OnActivity(GenerateToken, mock.Anything, mock.Anything).Return("mock-token", nil)

	// Setup delayed signal
	env.RegisterDelayedCallback(func() {
		s.Require().NotEmpty(generatedCode, "Code should have been captured from activity")
		s.T().Logf("Signaling workflow with code: %s", generatedCode)
		env.SignalWorkflow("SubmitCode", generatedCode)
	}, 1*time.Second)

	// Execute
	env.ExecuteWorkflow(LoginWorkflow, "test@example.com")

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	var result string
	err := env.GetWorkflowResult(&result)
	s.NoError(err)
	s.Equal("mock-token", result)
}

func (s *UnitTestSuite) TestLoginWorkflow_Timeout() {
	env := s.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(LoginWorkflow)

	// Mock SendMagicCode (we don't care about the code here, just that it sends)
	env.OnActivity(SendMagicCode, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Setup delayed callback to "do nothing" but verify timeout?
	// Actually, we just need to NOT signal.
	// We cannot easily "wait" in the test for 10 minutes without firing timers.
	// We explicitly verify that after 10+ minutes it fails.
	// But basic execution runs to completion (or deadlock if blocked).
	// TestSuite has automatic time skipping?
	// Yes, if we block, it skips to next timer.

	// So we just execute. The workflow will block on Selector.
	// The test environment should auto-fast-forward time to the timer firing
	// IF there are no other events.

	env.ExecuteWorkflow(LoginWorkflow, "test@example.com")

	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())
	s.Equal("workflow execution error (type: LoginWorkflow, workflowID: default-test-workflow-id, runID: default-test-run-id): Login timed out", env.GetWorkflowError().Error())
}
