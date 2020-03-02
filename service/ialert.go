package service

type ialert interface {
	Alert(title string, content string) bool
}
