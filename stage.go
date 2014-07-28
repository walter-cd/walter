package plumber

type Stage interface {
	Run() bool
}
