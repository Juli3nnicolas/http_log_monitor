package task

// Task describes a piece of work that can be executed on a delimited time-frame
type Task interface {
	// Init sets up the task. Think of it as a constructor.
	Init(...interface{}) error
	// BeforeRun sets some data before the task is run. Remember that
	// a task can be executed several times, accross several time frames.
	// This function and AfterRun is there to setup/cleanup the task's state
	// in case a long iterative work is planned
	BeforeRun() error
	// Run executes the task
	Run() error
	// AfterRun like BeforeRun except it is called after Run has been executed
	AfterRun() error
	// IsDone is true if the task's work has been completed
	IsDone() bool
	// Close shutdowns the task. Call Init to use it again
	Close() error
}
