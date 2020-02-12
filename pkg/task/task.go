package task

// Task describes a piece of computing
// To start it call Run
// To know whether it's over, call IsDone
type Task interface {
	Run() error
	IsDone() bool
}
