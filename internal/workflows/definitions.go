package workflows

const (
	TaskQueue = "QUEUE_TASK_QUEUE"

	// Signals
	SignalJoinQueue  = "JoinQueue"
	SignalLeaveQueue = "LeaveQueue"
	SignalCallNext   = "CallNext"
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

type CallNextSignal struct {
	CounterID string
}

type CallNextParams struct {
	BusinessID string
	UserID     string
	CounterID  string
	Status     string
}
