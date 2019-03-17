package configer

type Configer interface {
	Load(file string) error
	Get(name string) (value string, err error)
	MustGet(name string, defaultValue string) (value string)
}
