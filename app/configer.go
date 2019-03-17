package app

type ConfigLoader interface {
	Load(file string) error
}
