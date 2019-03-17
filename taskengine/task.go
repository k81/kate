package taskengine

// TaskFunc define the task func type
type TaskFunc func()

// Run implements the Task interface
func (f TaskFunc) Run() {
	f()
}

// Task define the task interface
type Task interface {
	Run()
}
