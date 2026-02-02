package workflows

const (
	TaskQueue = "QUEUE_TASK_QUEUE"

	// Signals
	SignalJoinQueue  = "JoinQueue"
	SignalLeaveQueue = "LeaveQueue"
	SignalExit       = "Exit" // Added for clean shutdown

	// Queries
	QueryGetState = "GetState"
)

type JoinQueueSignal struct {
	UserID string
}

type LeaveQueueSignal struct {
	UserID string
}
