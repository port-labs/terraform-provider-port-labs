package consts

const (
	Failure             = "FAILURE"
	Cancelled           = "CANCELLED"
	Completed           = "COMPLETED"
	Running             = "RUNNING"
	Pending             = "PENDING"
	Initializing        = "INITIALIZING"
	PendingCancellation = "PENDING_CANCELLATION"
)

func IsTerminalStatus(status string) bool {
	return status == Failure || status == Cancelled || status == Completed
}
