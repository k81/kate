package taskengine

type TaskFunc func()

func (f TaskFunc) Run() {
	f()
}

type Task interface {
	Run()
}
